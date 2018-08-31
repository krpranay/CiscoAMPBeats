Name: CiscoAMPBeats
Version: 1
Release: 0
Summary: Logrhythm CiscoAMPBeats

Group: Applications/System
License: (c) Copyright 2003-2012, LogRhythm, Inc.  All rights reserved.
URL: http://www.logrhythm.com

%description
Logrhythm CiscoAMPBeats

%files
%dir %attr(0700, root, root) %config /opt/logrhythm/
%dir %attr(0700, root, root) %config /opt/logrhythm/CiscoAMPBeats
%attr(0700, root, root) /opt/logrhythm/CiscoAMPBeats/CiscoAMPBeats
%attr(0600, root, root) /opt/logrhythm/CiscoAMPBeats/CiscoAMPBeats.yml
%attr(0600, root, root) /opt/logrhythm/CiscoAMPBeats/CiscoAMPEndPoint.ini
