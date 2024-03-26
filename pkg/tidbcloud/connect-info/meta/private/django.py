DATABASES = {
    'default': {
        'ENGINE': 'django_tidb',
        'NAME': '${database}',
        'USER': '${username}',
        'PASSWORD': '${password}',
        'HOST': '${host}',
        'PORT': ${port},
    },
}
