%define  debug_package %{nil}

Name:		gcredstash
Version:	0.3.5
Release:	1%{?dist}
Summary:	gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

Group:		Development/Tools
License:	Apache License, Version 2.0
URL:		https://github.com/kgaughan/gcredstash
Source0:	%{name}.tar.gz
# https://github.com/kgaughan/gcredstash/releases/download/v%{version}/gcredstash_%{version}.tar.gz

%description
gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

%prep
%setup -q -n src

%build
make VERSION=%{version}

%install
rm -rf %{buildroot}
mkdir -p %{buildroot}%{_sbindir}
install -m 700 gcredstash %{buildroot}%{_sbindir}

%files
%defattr(700,root,root,-)
%{_sbindir}/gcredstash
