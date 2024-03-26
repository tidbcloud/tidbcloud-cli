DATABASES = {
    'default': {
        'ENGINE': 'django_tidb',
        'NAME': '${database}',
        'USER': '${username}',
        'PASSWORD': '${password}',
        'HOST': '${host}',
        'PORT': ${port},
        'OPTIONS': {
            'ssl_mode': 'VERIFY_IDENTITY',
            'ssl': {'ca': '${ca_path}'}
        }
    },
}
