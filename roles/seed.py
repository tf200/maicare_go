import yaml
import psycopg2
from psycopg2.extras import DictCursor
import os
from dotenv import load_dotenv # type: ignore
import json

def load_yaml_file(file_path):
    with open(file_path, 'r') as file:
        return yaml.safe_load(file)

def connect_to_db():
    load_dotenv('../app.env')
    db_url = os.getenv('DB_SOURCE')
    return psycopg2.connect(db_url)

def insert_permissions(cursor, permissions):
    for permission in permissions:
        # Convert method list to string
        method_str = json.dumps(permission['method'])
        
        # Check if permission already exists
        cursor.execute("""
            SELECT id FROM permissions 
            WHERE name = %s AND resource = %s AND method = %s
        """, (permission['name'], permission['resource'], method_str))
        
        result = cursor.fetchone()
        
        if result is None:
            # Insert new permission
            cursor.execute("""
                INSERT INTO permissions (name, resource, method)
                VALUES (%s, %s, %s)
                RETURNING id
            """, (permission['name'], permission['resource'], method_str))
            
        print(f"Processed permission: {permission['name']}")

def insert_roles(cursor, roles):
    for role in roles:
        # Check if role exists
        cursor.execute("SELECT id FROM roles WHERE name = %s", (role['name'],))
        result = cursor.fetchone()
        
        role_id = None
        if result is None:
            # Insert new role
            cursor.execute("""
                INSERT INTO roles (id, name)
                VALUES (%s, %s)
                RETURNING id
            """, (role['id'], role['name']))
            role_id = cursor.fetchone()[0]
        else:
            role_id = result[0]
            
        print(f"Processed role: {role['name']}")
        
        # Handle permissions for this role
        if 'permissions' in role:
            for permission_name in role['permissions']:
                # Get permission id
                cursor.execute("SELECT id FROM permissions WHERE name = %s", (permission_name,))
                perm_result = cursor.fetchone()
                
                if perm_result:
                    permission_id = perm_result[0]
                    
                    # Check if role-permission mapping exists
                    cursor.execute("""
                        SELECT 1 FROM role_permissions 
                        WHERE role_id = %s AND permission_id = %s
                    """, (role_id, permission_id))
                    
                    if not cursor.fetchone():
                        # Insert new role-permission mapping
                        cursor.execute("""
                            INSERT INTO role_permissions (role_id, permission_id)
                            VALUES (%s, %s)
                        """, (role_id, permission_id))
                        print(f"Added permission {permission_name} to role {role['name']}")

def main():
    try:
        # Load YAML data
        data = load_yaml_file('rbac_config.yaml')
        
        # Connect to database
        conn = connect_to_db()
        cursor = conn.cursor(cursor_factory=DictCursor)
        
        try:
            # Begin transaction
            print("Starting database sync...")
            
            # Insert permissions
            if 'permissions' in data:
                print("\nProcessing permissions...")
                insert_permissions(cursor, data['permissions'])
            
            # Insert roles and their permissions
            if 'roles' in data:
                print("\nProcessing roles...")
                insert_roles(cursor, data['roles'])
            
            # Commit transaction
            conn.commit()
            print("\nDatabase sync completed successfully!")
            
        except Exception as e:
            conn.rollback()
            print(f"Error during database sync: {str(e)}")
            raise
        
        finally:
            cursor.close()
            conn.close()
            
    except Exception as e:
        print(f"Script failed: {str(e)}")
        raise

if __name__ == "__main__":
    main()