from setuptools import setup, find_packages

setup(
    name="cp-cli",
    version="0.1",
    packages=find_packages(),
    py_modules=["main"],
    install_requires=[
        "rich",
        "InquirerPy",
        "requests"
    ],
    entry_points={
        "console_scripts": [
            "cp-cli = main:main"
        ]
    },
    include_package_data=True
)