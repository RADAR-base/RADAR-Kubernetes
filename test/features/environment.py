def before_all(context):
    context.cache = {
        "management_portal_token": None,
        "armt_source_type_json": None,
        "organization_json": None,
        "project_json": None,
        "armt_project_source_json": None,
        "test_subject_id": None,
        "secrets": None,
        "armt_meta_token": None,
        "armt_refresh_token": None,
        "armt_access_token": None,
        "rest_auth_registration_json": None,
        "fitbit_user_json": None,
    }
    context.state = {
        "database": {},
        "storage": {},
    }
