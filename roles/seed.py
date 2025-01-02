import yaml
import psycopg2
from psycopg2 import sql
import os
from dotenv import load_dotenv # type: ignore
from pathlib import Path



env_path = Path(__file__).resolve().parent.parent / "app.env"

# Load environment variables from app.env
load_dotenv(env_path)

# Get the database connection string from the environment variable
DB_SOURCE = os.getenv("DB_SOURCE")
if not DB_SOURCE:
    raise ValueError("DB_SOURCE is not set in environment variables.")

# Path to the YAML file
YAML_PATH = os.path.join(os.path.dirname(__file__), "rbac_config.yaml")

# Connect to the database
def connect_db():
    conn = psycopg2.connect(DB_SOURCE)
    return conn

# Seed resources and permissions
def seed_resources_permissions(conn, resources):
    with conn.cursor() as cursor:
        for resource in resources:
            for method in resource["methods"]:
                cursor.execute(
                    sql.SQL("""
                        INSERT INTO Permissions (name, resource)
                        VALUES (%s, %s)
                        ON CONFLICT DO NOTHING
                    """), (method["name"], resource["path"])
                )
    conn.commit()

# Seed roles and role_permissions
def seed_roles_permissions(conn, roles):
    with conn.cursor() as cursor:
        for role in roles:
            # Insert role
            cursor.execute(
                sql.SQL("""
                    INSERT INTO Roles (name)
                    VALUES (%s)
                    ON CONFLICT (name) DO UPDATE
                    SET name = EXCLUDED.name
                    RETURNING id
                """), (role["name"],)
            )
            role_id = cursor.fetchone()[0]

            # Insert role_permissions
            for permission in role["permissions"]:
                for method in permission["methods"]:
                    cursor.execute(
                        sql.SQL("""
                            SELECT id FROM Permissions
                            WHERE name = %s AND resource = %s
                        """), (method, permission["resource"])
                    )
                    permission_id = cursor.fetchone()[0]
                    cursor.execute(
                        sql.SQL("""
                            INSERT INTO Role_Permissions (role_id, permission_id)
                            VALUES (%s, %s)
                            ON CONFLICT DO NOTHING
                        """), (role_id, permission_id)
                    )
    conn.commit()

# Main function
def main():
    # Read the YAML file
    with open(YAML_PATH, "r") as file:
        config = yaml.safe_load(file)

    # Connect to the database
    conn = connect_db()

    # Seed resources and permissions
    seed_resources_permissions(conn, config["resources"])

    # Seed roles and role_permissions
    seed_roles_permissions(conn, config["roles"])

    # Close the connection
    conn.close()
    print("Database seeding completed successfully!")

if __name__ == "__main__":
    main()