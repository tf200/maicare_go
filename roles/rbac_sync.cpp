#include <iostream>
#include <fstream>
#include <sstream>
#include <cstdlib>
#include <string>
#include <stdexcept>
#include <vector>
#include <algorithm>
#include <cctype>
#include <locale>

// Include yaml-cpp for YAML parsing
#include <yaml-cpp/yaml.h>

// Include libpqxx for PostgreSQL connectivity
#include <pqxx/pqxx>

// Helper functions to trim whitespace
static inline void ltrim(std::string &s) {
    s.erase(s.begin(), std::find_if(s.begin(), s.end(),
        [](unsigned char ch) { return !std::isspace(ch); }));
}

static inline void rtrim(std::string &s) {
    s.erase(std::find_if(s.rbegin(), s.rend(),
        [](unsigned char ch) { return !std::isspace(ch); }).base(), s.end());
}

static inline void trim(std::string &s) {
    ltrim(s);
    rtrim(s);
}

// Helper function to load a specific key's value from the env file.
std::string load_env_variable(const std::string &filename, const std::string &key) {
    std::ifstream file(filename);
    if (!file.is_open()) {
        throw std::runtime_error("Could not open env file: " + filename);
    }
    
    std::string line;
    while (std::getline(file, line)) {
        // Remove any whitespace from beginning and end.
        trim(line);
        // Skip empty lines or comments.
        if (line.empty() || line[0] == '#') {
            continue;
        }
        // Look for key=value
        size_t pos = line.find('=');
        if (pos == std::string::npos) {
            continue;
        }
        std::string var = line.substr(0, pos);
        std::string value = line.substr(pos + 1);
        trim(var);
        trim(value);
        if (var == key) {
            return value;
        }
    }
    throw std::runtime_error("Key " + key + " not found in " + filename);
}

// Helper function to serialize a YAML sequence as a JSON‑like string
std::string serializeMethod(const YAML::Node &node) {
    if (!node.IsSequence()) {
        return "\"" + node.as<std::string>() + "\"";
    }
    std::string result = "[";
    bool first = true;
    for (const auto &item : node) {
        if (!first) {
            result += ", ";
        }
        result += "\"" + item.as<std::string>() + "\"";
        first = false;
    }
    result += "]";
    return result;
}

// Insert permissions into the database
void insert_permissions(pqxx::work &txn, const YAML::Node &permissions) {
    for (const auto &perm : permissions) {
        std::string name     = perm["name"].as<std::string>();
        std::string resource = perm["resource"].as<std::string>();
        std::string method   = serializeMethod(perm["method"]);

        std::string selectQuery =
            "SELECT id FROM permissions "
            "WHERE name = "     + txn.quote(name) +
            "  AND resource = " + txn.quote(resource) +
            "  AND method = "   + txn.quote(method) + ";";

        pqxx::result res = txn.exec(selectQuery);

        if (res.empty()) {
            std::string insertQuery =
                "INSERT INTO permissions (name, resource, method) VALUES (" +
                txn.quote(name) + ", " +
                txn.quote(resource) + ", " +
                txn.quote(method) + ");";
            txn.exec(insertQuery);
            std::cout << "Inserted permission: " << name << std::endl;
        } else {
            std::cout << "Permission already exists: " << name << std::endl;
        }
    }
}


void insert_roles(pqxx::work &txn, const YAML::Node &roles) {
    for (const auto &role : roles) {
        std::string name        = role["name"].as<std::string>();
        int         id          = role["id"].as<int>();
        std::string description = role["description"]
                                      ? role["description"].as<std::string>()
                                      : "";

        // Upsert role row
        std::string upsertRole =
            "INSERT INTO roles ( name) VALUES (" +
            txn.quote(name) + ") " +
            "ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;";
        txn.exec(upsertRole);
        std::cout << "Processed role: " << name << std::endl;

        // Sync role_permissions
        if (role["permissions"]) {
            for (const auto &permNameNode : role["permissions"]) {
                std::string permName = permNameNode.as<std::string>();

                // Get permission id
                pqxx::result permRes =
                    txn.exec("SELECT id FROM permissions WHERE name = " +
                             txn.quote(permName) + ";");
                if (permRes.empty()) {
                    std::cerr << "Warning: permission '" << permName
                              << "' not found, skipping.\n";
                    continue;
                }
                int permId = permRes[0]["id"].as<int>();

                // Idempotent insert into role_permissions
                std::string insertMap =
                    "INSERT INTO role_permissions (role_id, permission_id) VALUES (" +
                    txn.quote(id) + ", " + txn.quote(permId) + ") " +
                    "ON CONFLICT (role_id, permission_id) DO NOTHING;";
                txn.exec(insertMap);
            }
        }
    }
}

int main() {
    try {
        // Load YAML configuration from file
        YAML::Node config = YAML::LoadFile("rbac_config.yaml");

        // Read DB_SOURCE from app.env (adjust the path as needed)
        std::string db_source = load_env_variable("../app.env", "DB_SOURCE");
        std::cout << "DB_SOURCE: " << db_source << std::endl;

        // Connect to the PostgreSQL database
        pqxx::connection conn(db_source);
        pqxx::work txn(conn);  // Begin transaction

        std::cout << "Starting database sync..." << std::endl;

        // Process permissions
        if (config["permissions"]) {
            std::cout << "\nProcessing permissions..." << std::endl;
            insert_permissions(txn, config["permissions"]);
        }

        // Process roles and role-permission mappings
        if (config["roles"]) {
            std::cout << "\nProcessing roles..." << std::endl;
            insert_roles(txn, config["roles"]);
        }

        // Commit transaction
        txn.commit();
        std::cout << "\nDatabase sync completed successfully!" << std::endl;
    }
    catch (const std::exception &e) {
        std::cerr << "Script failed: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}
