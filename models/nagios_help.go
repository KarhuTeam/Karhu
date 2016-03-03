package models

var NagiosHelps = map[string]string{
	"check_apt": `check_apt v1.5 (nagios-plugins 1.5)
Copyright (c) 2006-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks for software updates on systems that use
package management systems based on the apt-get(8) command
found in Debian GNU/Linux


Usage:
check_apt [[-d|-u|-U]opts] [-n] [-t timeout]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -U, --upgrade=OPTS
    [Default] Perform an upgrade.  If an optional OPTS argument is provided,
    apt-get will be run with these command line options instead of the
    default (-o 'Debug::NoLocking=true' -s -qq).
    Note that you may be required to have root privileges if you do not use
    the default options.
 -d, --dist-upgrade=OPTS
    Perform a dist-upgrade instead of normal upgrade. Like with -U OPTS
    can be provided to override the default options.
  -n, --no-upgrade
    Do not run the upgrade.  Probably not useful (without -u at least).
 -i, --include=REGEXP
    Include only packages matching REGEXP.  Can be specified multiple times
    the values will be combined together.  Any packages matching this list
    cause the plugin to return WARNING status.  Others will be ignored.
    Default is to include all packages.
 -e, --exclude=REGEXP
    Exclude packages matching REGEXP from the list of packages that would
    otherwise be included.  Can be specified multiple times; the values
    will be combined together.  Default is to exclude no packages.
 -c, --critical=REGEXP
    If the full package information of any of the upgradable packages match
    this REGEXP, the plugin will return CRITICAL status.  Can be specified
    multiple times like above.  Default is a regexp matching security
    upgrades for Debian and Ubuntu:
    	^[^\(]*\(.* (Debian-Security:|Ubuntu:[^/]*/[^-]*-security)
    Note that the package must first match the include list before its
    information is compared against the critical list.

The following options require root privileges and should be used with care:

 -u, --update=OPTS
    First perform an 'apt-get update'.  An optional OPTS parameter overrides
    the default options.  Note: you may also need to adjust the global
    timeout (with -t) to prevent the plugin from timing out if apt-get
    upgrade is expected to take longer than the default timeout.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_breeze": `check_breeze v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2000 Jeffrey Blank/Karl DeBisschop

This plugin reports the signal strength of a Breezecom wireless equipment

Usage: check_breeze -H <host> [-C community] -w <warn> -c <crit>

-H, --hostname=HOST
   Name or IP address of host to check
-C, --community=community
   SNMPv1 community (default public)
-w, --warning=INTEGER
   Percentage strength below which a WARNING status will result
-c, --critical=INTEGER
   Percentage strength below which a CRITICAL status will result

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_by_ssh": `check_by_ssh v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Karl DeBisschop <kdebisschop@users.sourceforge.net>
Copyright (c) 2000-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin uses SSH to execute commands on a remote host

Usage:
 check_by_ssh -H <host> -C <command> [-fqv] [-1|-2] [-4|-6]
       [-S [lines]] [-E [lines]] [-t timeout] [-i identity]
       [-l user] [-n name] [-s servicelist] [-O outputfile]
       [-p port] [-o ssh-option] [-F configfile]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -1, --proto1
    tell ssh to use Protocol 1 [optional]
 -2, --proto2
    tell ssh to use Protocol 2 [optional]
 -S, --skip-stdout[=n]
    Ignore all or (if specified) first n lines on STDOUT [optional]
 -E, --skip-stderr[=n]
    Ignore all or (if specified) first n lines on STDERR [optional]
 -f
    tells ssh to fork rather than create a tty [optional]. This will always return OK if ssh is executed
 -C, --command='COMMAND STRING'
    command to execute on the remote machine
 -l, --logname=USERNAME
    SSH user name on remote host [optional]
 -i, --identity=KEYFILE
    identity of an authorized key [optional]
 -O, --output=FILE
    external command file for nagios [optional]
 -s, --services=LIST
    list of nagios service names, separated by ':' [optional]
 -n, --name=NAME
    short name of host in nagios configuration [optional]
 -o, --ssh-option=OPTION
    Call ssh with '-o OPTION' (may be used multiple times) [optional]
 -F, --configfile
    Tell ssh to use this configfile [optional]
 -q, --quiet
    Tell ssh to suppress warning and diagnostic messages [optional]
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

 The most common mode of use is to refer to a local identity file with
 the '-i' option. In this mode, the identity pair should have a null
 passphrase and the public key should be listed in the authorized_keys
 file of the remote host. Usually the key will be restricted to running
 only one command on the remote server. If the remote SSH server tracks
 invocation arguments, the one remote program may be an agent that can
 execute additional commands as proxy

 To use passive mode, provide multiple '-C' options, and provide
 all of -O, -s, and -n options (servicelist order must match '-C'options)

Examples:
 $ check_by_ssh -H localhost -n lh -s c1:c2:c3 -C uptime -C uptime -C uptime -O /tmp/foo
 $ cat /tmp/foo
 [1080933700] PROCESS_SERVICE_CHECK_RESULT;flint;c1;0; up 2 days
 [1080933700] PROCESS_SERVICE_CHECK_RESULT;flint;c2;0; up 2 days
 [1080933700] PROCESS_SERVICE_CHECK_RESULT;flint;c3;0; up 2 days

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_clamd": `check_clamd v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests CLAMD connections with the specified host (or unix socket).

Usage:
check_clamd -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_cluster": `Could not parse arguments
Usage:
 check_cluster (-s | -h) -d val1[,val2,...,valn] [-l label]
[-w threshold] [-c threshold] [-v] [--help]`,
	"check_dbi": `check_dbi v1.5 (nagios-plugins 1.5)
Copyright (c) 2011 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This program connects to an (SQL) database using DBI and checks the
specified metric against threshold levels. The default metric is
the result of the specified query.


Usage:
check_dbi -d <DBI driver> [-o <DBI driver option> [...]] [-q <query>]
 [-H <host>] [-c <critical range>] [-w <warning range>] [-m <metric>]
 [-e <string>] [-r|-R <regex>]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.

 -d, --driver=STRING
    DBI driver to use
 -o, --option=STRING
    DBI driver options
 -q, --query=STRING
    query to execute

 -w, --warning=RANGE
    Warning range (format: start:end). Alert if outside this range
 -c, --critical=RANGE
    Critical range
 -e, --expect=STRING
    String to expect as query result
    Do not mix with -w, -c, -r, or -R!
 -r, --regex=REGEX
    Extended POSIX regular expression to check query result against
    Do not mix with -w, -c, -e, or -R!
 -R, --regexi=REGEX
    Case-insensitive extended POSIX regex to check query result against
    Do not mix with -w, -c, -e, or -r!
 -m, --metric=METRIC
    Metric to check thresholds against. Available metrics:
    CONN_TIME    - time used for setting up the database connection
    QUERY_RESULT - result (first column of first row) of the query
    QUERY_TIME   - time used to execute the query
                   (ignore the query result)

 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

 A DBI driver (-d option) is required. If the specified metric operates
 on a query, one has to be specified (-q option).

 This plugin connects to an (SQL) database using libdbi and, optionally,
 executes the specified query. The first column of the first row of the
 result will be parsed and, in QUERY_RESULT mode, compared with the
 warning and critical ranges. The result from the query has to be numeric
 (strings representing numbers are fine).

 The number and type of required DBI driver options depends on the actual
 driver. See its documentation at http://libdbi-drivers.sourceforge.net/
 for details.

 Examples:
  check_dbi -d pgsql -o username=postgres -m QUERY_RESULT \
    -q 'SELECT COUNT(*) FROM pg_stat_activity' -w 5 -c 10
  Warning if more than five connections; critical if more than ten.

  check_dbi -d mysql -H localhost -o username=user -o password=secret \
    -q 'SELECT COUNT(*) FROM logged_in_users -w 5:20 -c 0:50
  Warning if less than 5 or more than 20 users are logged in; critical
  if more than 50 users.

  check_dbi -d firebird -o username=user -o password=secret -o dbname=foo \
    -m CONN_TIME -w 0.5 -c 2
  Warning if connecting to the database takes more than half of a second;
  critical if it takes more than 2 seconds.

  check_dbi -d mysql -H localhost -o username=user \
    -q 'SELECT concat(@@version, " ", @@version_comment)' \
    -r '^5\.[01].*MySQL Enterprise Server'
  Critical if the database server is not a MySQL enterprise server in either
  version 5.0.x or 5.1.x.

  check_dbi -d pgsql -u username=user -m SERVER_VERSION \
    -w 090000:090099 -c 090000:090199
  Warn if the PostgreSQL server version is not 9.0.x; critical if the version
  is less than 9.x or higher than 9.1.x.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_dhcp": `check_dhcp v1.5 (nagios-plugins 1.5)
Copyright (c) 2001-2004 Ethan Galstad (nagios@nagios.org)
Copyright (c) 2001-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests the availability of DHCP servers on a network.


Usage:
 check_dhcp [-v] [-u] [-s serverip] [-r requestedip] [-t timeout]
                  [-i interface] [-m mac]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)
 -s, --serverip=IPADDRESS
    IP address of DHCP server that we must hear from
 -r, --requestedip=IPADDRESS
    IP address that should be offered by at least one DHCP server
 -t, --timeout=INTEGER
    Seconds to wait for DHCPOFFER before timeout occurs
 -i, --interface=STRING
    Interface to to use for listening (i.e. eth0)
 -m, --mac=STRING
    MAC address to use in the DHCP request
 -u, --unicast
    Unicast testing: mimic a DHCP relay, requires -s

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_dig": `check_dig v1.5 (nagios-plugins 1.5)
Copyright (c) 2000 Karl DeBisschop <kdebisschop@users.sourceforge.net>
Copyright (c) 2002-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin test the DNS service on the specified host using dig

Usage:
check_dig -l <query_address> [-H <host>] [-p <server port>]
 [-T <query type>] [-w <warning interval>] [-c <critical interval>]
 [-t <timeout>] [-a <expected answer address>] [-v]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 53)
 -4, --use-ipv4
    Force dig to only use IPv4 query transport
 -6, --use-ipv6
    Force dig to only use IPv6 query transport
 -l, --query_address=STRING
    Machine name to lookup
 -T, --record_type=STRING
    Record type to lookup (default: A)
 -a, --expected_address=STRING
    An address expected to be in the answer section. If not set, uses whatever
    was in -l
 -A, --dig-arguments=STRING
    Pass STRING as argument(s) to dig
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Examples:
 check_dig -H DNSSERVER -l www.example.com -A "+tcp"
 This will send a tcp query to DNSSERVER for www.example.com

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_disk": `check_disk v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks the amount of used disk space on a mounted file system
and generates an alert if free space is less than one of the threshold values


Usage:
 check_disk -w limit -c limit [-W limit] [-K limit] {-p path | -x device}
[-C] [-E] [-e] [-f] [-g group ] [-k] [-l] [-M] [-m] [-R path ] [-r path ]
[-t timeout] [-u unit] [-v] [-X type] [-N type]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -w, --warning=INTEGER
    Exit with WARNING status if less than INTEGER units of disk are free
 -w, --warning=PERCENT%
    Exit with WARNING status if less than PERCENT of disk space is free
 -c, --critical=INTEGER
    Exit with CRITICAL status if less than INTEGER units of disk are free
 -c, --critical=PERCENT%
    Exit with CRITICAL status if less than PERCENT of disk space is free
 -W, --iwarning=PERCENT%
    Exit with WARNING status if less than PERCENT of inode space is free
 -K, --icritical=PERCENT%
    Exit with CRITICAL status if less than PERCENT of inode space is free
 -p, --path=PATH, --partition=PARTITION
    Path or partition (may be repeated)
 -x, --exclude_device=PATH <STRING>
    Ignore device (only works if -p unspecified)
 -C, --clear
    Clear thresholds
 -E, --exact-match
    For paths or partitions specified with -p, only check for exact paths
 -e, --errors-only
    Display only devices/mountpoints with errors
 -f, --freespace-ignore-reserved
    Don't account root-reserved blocks into freespace in perfdata
 -g, --group=NAME
    Group paths. Thresholds apply to (free-)space of all partitions together
 -k, --kilobytes
    Same as '--units kB'
 -l, --local
    Only check local filesystems
 -L, --stat-remote-fs
    Only check local filesystems against thresholds. Yet call stat on remote filesystems
    to test if they are accessible (e.g. to detect Stale NFS Handles)
 -M, --mountpoint
    Display the mountpoint instead of the partition
 -m, --megabytes
    Same as '--units MB'
 -A, --all
    Explicitly select all paths. This is equivalent to -R '.*'
 -R, --eregi-path=PATH, --eregi-partition=PARTITION
    Case insensitive regular expression for path/partition (may be repeated)
 -r, --ereg-path=PATH, --ereg-partition=PARTITION
    Regular expression for path or partition (may be repeated)
 -I, --ignore-eregi-path=PATH, --ignore-eregi-partition=PARTITION
    Regular expression to ignore selected path/partition (case insensitive) (may be repeated)
 -i, --ignore-ereg-path=PATH, --ignore-ereg-partition=PARTITION
    Regular expression to ignore selected path or partition (may be repeated)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -u, --units=STRING
    Choose bytes, kB, MB, GB, TB (default: MB)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)
 -X, --exclude-type=TYPE
    Ignore all filesystems of indicated type (may be repeated)
 -N, --include-type=TYPE
    Check only filesystems of indicated type (may be repeated)

Examples:
 check_disk -w 10% -c 5% -p /tmp -p /var -C -w 100000 -c 50000 -p /
    Checks /tmp and /var at 10% and 5%, and / at 100MB and 50MB
 check_disk -w 100 -c 50 -C -w 1000 -c 500 -g sidDATA -r '^/oracle/SID/data.*$'
    Checks all filesystems not matching -r at 100M and 50M. The fs matching the -r regex
    are grouped which means the freespace thresholds are applied to all disks together
 check_disk -w 100 -c 50 -C -w 1000 -c 500 -p /foo -C -w 5% -c 3% -p /bar
    Checks /foo for 1000M/500M and /bar for 5/3%. All remaining volumes use 100M/50M

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_disk_smb": `check_disk_smb v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2000 Michael Anthon/Karl DeBisschop

Perl Check SMB Disk plugin for Nagios

Usage: check_disk_smb -H <host> -s <share> -u <user> -p <password>
      -w <warn> -c <crit> [-W <workgroup>] [-P <port>] [-a <IP>]

-H, --hostname=HOST
   NetBIOS name of the server
-s, --share=STRING
   Share name to be tested
-W, --workgroup=STRING
   Workgroup or Domain used (Defaults to "WORKGROUP")
-a, --address=IP
   IP-address of HOST (only necessary if HOST is in another network)
-u, --user=STRING
   Username to log in to server. (Defaults to "guest")
-p, --password=STRING
   Password to log in to server. (Defaults to an empty password)
-w, --warning=INTEGER or INTEGER[kMG]
   Percent of used space at which a warning will be generated (Default: 85%)

-c, --critical=INTEGER or INTEGER[kMG]
   Percent of used space at which a critical will be generated (Defaults: 95%)
-P, --port=INTEGER
   Port to be used to connect to. Some Windows boxes use 139, others 445 (Defaults to smbclient default)

   If thresholds are followed by either a k, M, or G then check to see if that
   much disk space is available (kilobytes, Megabytes, Gigabytes)

   Warning percentage should be less than critical
   Warning (remaining) disk space should be greater than critical.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_dns": `check_dns v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 2000-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin uses the nslookup program to obtain the IP address for the given host/domain query.
An optional DNS server to use may be specified.
If no DNS server is specified, the default server(s) specified in /etc/resolv.conf will be used.


Usage:
check_dns -H host [-s server] [-a expected-address] [-A] [-t timeout] [-w warn] [-c crit]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=HOST
    The name or address you want to query
 -s, --server=HOST
    Optional DNS server you want to use for the lookup
 -a, --expected-address=IP-ADDRESS|HOST
    Optional IP-ADDRESS you expect the DNS server to return. HOST must end with
    a dot (.). This option can be repeated multiple times (Returns OK if any
    value match). If multiple addresses are returned at once, you have to match
    the whole string of addresses separated with commas (sorted alphabetically).
 -A, --expect-authority
    Optionally expect the DNS server to be authoritative for the lookup
 -w, --warning=seconds
    Return warning if elapsed time exceeds value. Default off
 -c, --critical=seconds
    Return critical if elapsed time exceeds value. Default off
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_dummy": `check_dummy v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin will simply return the state corresponding to the numeric value
of the <state> argument with optional text


Usage:
 check_dummy <integer state> [optional text]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_file_age": `check_file_age v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2003 Steven Grimm

Usage:
  check_file_age [-w <secs>] [-c <secs>] [-W <size>] [-C <size>] -f <file>
  check_file_age [-h | --help]
  check_file_age [-V | --version]

  <secs>  File must be no more than this many seconds old (default: warn 240 secs, crit 600)
  <size>  File must be at least this many bytes long (default: crit 0 bytes)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_flexlm": `check_flexlm v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2000 Ernst-Dieter Martin/Karl DeBisschop

Check available flexlm license managers

Usage:
   check_flexlm -F <filename> [-v] [-t] [-V] [-h]
   check_flexlm --help
   check_flexlm --version

-F, --filename=FILE
   Name of license file (usually "license.dat")
-v, --verbose
   Print some extra debugging information (not advised for normal operation)
-t, --timeout
   Plugin time out in seconds (default = 15 )
-V, --version
   Show version and license information
-h, --help
   Show this help screen

Flexlm license managers usually run as a single server or three servers and a
quorum is needed.  The plugin return OK if 1 (single) or 3 (triple) servers
are running, CRITICAL if 1(single) or 3 (triple) servers are down, and WARNING
if 1 or 2 of 3 servers are running

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_ftp": `check_ftp v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests FTP connections with the specified host (or unix socket).

Usage:
check_ftp -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_host": `Copyright (c) 2005 Andreas Ericsson <ae@op5.se>
Copyright (c) 2005-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>



Usage:
 check_host [options] [-H] host1 host2 hostN

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H
    specify a target
 -w
    warning threshold (currently 1000,000ms,100%)
 -c
    critical threshold (currently 1000,000ms,100%)
 -s
    specify a source IP address or device name
 -n
    number of packets to send (currently 5)
 -i
    max packet interval (currently 1000,000ms)
 -I
    max target interval (currently 0,000ms)
 -m
    number of alive hosts required for success
 -l
    TTL on outgoing packets (currently 0)
 -t
    timeout value (seconds, currently  10)
 -b
    Number of icmp data bytes to send
    Packet size will be data bytes + icmp header (currently 68 + 8)
 -v
    verbose

Notes:
 The -H switch is optional. Naming a host (or several) to check is not.

 Threshold format for -w and -c is 200.25,60% for 200.25 msec RTA and 60%
 packet loss.  The default values should work well for most users.
 You can specify different RTA factors using the standardized abbreviations
 us (microseconds), ms (milliseconds, default) or just plain s for seconds.

 The -v switch can be specified several times for increased verbosity.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_hpjd": `check_hpjd v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests the STATUS of an HP printer with a JetDirect card.
Net-snmp must be installed on the computer running the plugin.


Usage:
check_hpjd -H host [-C community]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -C, --community=STRING
    The SNMP community name (default=public)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_http": `check_http v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2013 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests the HTTP service on the specified host. It can test
normal (http) and secure (https) servers, follow redirects, search for
strings and regular expressions, check connection times, and report on
certificate expiration times.


Usage:
 check_http -H <vhost> | -I <IP-address> [-u <uri>] [-p <port>]
       [-J <client certificate file>] [-K <private key>]
       [-w <warn time>] [-c <critical time>] [-t <timeout>] [-L] [-E] [-a auth]
       [-b proxy_auth] [-f <ok|warning|critcal|follow|sticky|stickyport>]
       [-e <expect>] [-d string] [-s string] [-l] [-r <regex> | -R <case-insensitive regex>]
       [-P string] [-m <min_pg_size>:<max_pg_size>] [-4|-6] [-N] [-M <age>]
       [-A string] [-k string] [-S <version>] [--sni] [-C <warn_age>[,<crit_age>]]
       [-T <content-type>] [-j method]
NOTE: One or both of -H and -I must be specified

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name argument for servers using host headers (virtual host)
    Append a port to include it in the header (eg: example.com:5000)
 -I, --IP-address=ADDRESS
    IP address or name (use numeric address if possible to bypass DNS lookup).
 -p, --port=INTEGER
    Port number (default: 80)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -S, --ssl=VERSION
    Connect via SSL. Port defaults to 443. VERSION is optional, and prevents
    auto-negotiation (1 = TLSv1, 2 = SSLv2, 3 = SSLv3).
 --sni
    Enable SSL/TLS hostname extension support (SNI)
 -C, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid. Port defaults to 443
    (when this option is used the URL is not checked.)
 -J, --client-cert=FILE
   Name of file that contains the client certificate (PEM format)
   to be used in establishing the SSL session
 -K, --private-key=FILE
   Name of file containing the private key (PEM format)
   matching the client certificate
 -e, --expect=STRING
    Comma-delimited list of strings, at least one of them is expected in
    the first (status) line of the server response (default: HTTP/1.)
    If specified skips all other status line logic (ex: 3xx, 4xx, 5xx processing)
 -d, --header-string=STRING
    String to expect in the response headers
 -s, --string=STRING
    String to expect in the content
 -u, --url=PATH
    URL to GET or POST (default: /)
 -P, --post=STRING
    URL encoded http POST data
 -j, --method=STRING  (for example: HEAD, OPTIONS, TRACE, PUT, DELETE)
    Set HTTP method.
 -N, --no-body
    Don't wait for document body: stop reading after headers.
    (Note that this still does an HTTP GET or POST, not a HEAD.)
 -M, --max-age=SECONDS
    Warn if document is more than SECONDS old. the number can also be of
    the form "10m" for minutes, "10h" for hours, or "10d" for days.
 -T, --content-type=STRING
    specify Content-Type header media type when POSTing

 -l, --linespan
    Allow regex to span newlines (must precede -r or -R)
 -r, --regex, --ereg=STRING
    Search page for regex STRING
 -R, --eregi=STRING
    Search page for case-insensitive regex STRING
 --invert-regex
    Return CRITICAL if found, OK if not

 -a, --authorization=AUTH_PAIR
    Username:password on sites with basic authentication
 -b, --proxy-authorization=AUTH_PAIR
    Username:password on proxy-servers with basic authentication
 -A, --useragent=STRING
    String to be sent in http header as "User Agent"
 -k, --header=STRING
    Any other tags to be sent in http header. Use multiple times for additional headers
 -E, --extended-perfdata
    Print additional performance data
 -L, --link
    Wrap output in HTML link (obsoleted by urlize)
 -f, --onredirect=<ok|warning|critical|follow|sticky|stickyport>
    How to handle redirected pages. sticky is like follow but stick to the
    specified IP address. stickyport also ensures port stays the same.
 -m, --pagesize=INTEGER<:INTEGER>
    Minimum page size required (bytes) : Maximum page size required (bytes)
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Notes:
 This plugin will attempt to open an HTTP connection with the host.
 Successful connects return STATE_OK, refusals and timeouts return STATE_CRITICAL
 other errors return STATE_UNKNOWN.  Successful connects, but incorrect reponse
 messages from the host result in STATE_WARNING return values.  If you are
 checking a virtual server that uses 'host headers' you must supply the FQDN
 (fully qualified domain name) as the [host_name] argument.

 This plugin can also check whether an SSL enabled web server is able to
 serve content (optionally within a specified time) or whether the X509
 certificate is still valid for the specified number of days.

 Please note that this plugin does not check if the presented server
 certificate matches the hostname of the server, or if the certificate
 has a valid chain of trust to one of the locally installed CAs.

Examples:
 CHECK CONTENT: check_http -w 5 -c 10 --ssl -H www.verisign.com

 When the 'www.verisign.com' server returns its content within 5 seconds,
 a STATE_OK will be returned. When the server returns its content but exceeds
 the 5-second threshold, a STATE_WARNING will be returned. When an error occurs,
 a STATE_CRITICAL will be returned.

 CHECK CERTIFICATE: check_http -H www.verisign.com -C 14

 When the certificate of 'www.verisign.com' is valid for more than 14 days,
 a STATE_OK is returned. When the certificate is still valid, but for less than
 14 days, a STATE_WARNING is returned. A STATE_CRITICAL will be returned when
 the certificate is expired.

 CHECK CERTIFICATE: check_http -H www.verisign.com -C 30,14

 When the certificate of 'www.verisign.com' is valid for more than 30 days,
 a STATE_OK is returned. When the certificate is still valid, but for less than
 30 days, but more than 14 days, a STATE_WARNING is returned.
 A STATE_CRITICAL will be returned when certificate expires in less than 14 days

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_icmp": `Copyright (c) 2005 Andreas Ericsson <ae@op5.se>
Copyright (c) 2005-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>



Usage:
 check_icmp [options] [-H] host1 host2 hostN

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H
    specify a target
 -w
    warning threshold (currently 200,000ms,40%)
 -c
    critical threshold (currently 500,000ms,80%)
 -s
    specify a source IP address or device name
 -n
    number of packets to send (currently 5)
 -i
    max packet interval (currently 80,000ms)
 -I
    max target interval (currently 0,000ms)
 -m
    number of alive hosts required for success
 -l
    TTL on outgoing packets (currently 0)
 -t
    timeout value (seconds, currently  10)
 -b
    Number of icmp data bytes to send
    Packet size will be data bytes + icmp header (currently 68 + 8)
 -v
    verbose

Notes:
 The -H switch is optional. Naming a host (or several) to check is not.

 Threshold format for -w and -c is 200.25,60% for 200.25 msec RTA and 60%
 packet loss.  The default values should work well for most users.
 You can specify different RTA factors using the standardized abbreviations
 us (microseconds), ms (milliseconds, default) or just plain s for seconds.

 The -v switch can be specified several times for increased verbosity.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ide_smart": `check_ide_smart v1.5 (nagios-plugins 1.5)
Nagios feature - 1999 Robert Dale <rdale@digital-mission.com>
(C) 1999 Ragnar Hojland Espinosa <ragnar@lightside.dhis.org>
Copyright (c) 1998-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks a local hard drive with the (Linux specific) SMART interface [http://smartlinux.sourceforge.net/smart/index.php].

Usage:
check_ide_smart [-d <device>] [-i <immediate>] [-q quiet] [-1 <auto-on>] [-O <auto-off>] [-n <nagios>]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -d, --device=DEVICE
    Select device DEVICE
    Note: if the device is selected with this option, _no_ other options are accepted
 -i, --immediate
    Perform immediately offline tests
 -q, --quiet-check
    Returns the number of failed tests
 -1, --auto-on
    Turn on automatic offline tests
 -0, --auto-off
    Turn off automatic offline tests
 -n, --nagios
    Output suitable for Nagios

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ifoperstatus": `check_ifoperstatus v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.

usage:
check_ifoperstatus -k <IF_KEY> -H <HOSTNAME> [-C <community>]
Copyright (C) 2000 Christoph Kron
check_ifoperstatus.pl comes with ABSOLUTELY NO WARRANTY
This programm is licensed under the terms of the GNU General Public License
(check source code for details)


check_ifoperstatus plugin for Nagios monitors operational
status of a particular network interface on the target host

Usage:
   -H (--hostname)   Hostname to query - (required)
   -C (--community)  SNMP read community (defaults to public,
                     used with SNMP v1 and v2c
   -v (--snmp_version)  1 for SNMP v1 (default)
                        2 for SNMP v2c
                        SNMP v2c will use get_bulk for less overhead
                        if monitoring with -d
   -L (--seclevel)   choice of "noAuthNoPriv", "authNoPriv", or	"authPriv"
   -U (--secname)    username for SNMPv3 context
   -c (--context)    SNMPv3 context name (default is empty string)
   -A (--authpass)   authentication password (cleartext ascii or localized key
                     in hex with 0x prefix generated by using "snmpkey" utility
                     auth password and authEngineID
   -a (--authproto)  Authentication protocol (MD5 or SHA1)
   -X (--privpass)   privacy password (cleartext ascii or localized key
                     in hex with 0x prefix generated by using "snmpkey" utility
                     privacy password and authEngineID
   -P (--privproto)  privacy protocol (DES or AES; default: DES)
   -k (--key)        SNMP IfIndex value
   -d (--descr)      SNMP ifDescr value
   -T (--type)       SNMP ifType integer value (see http://www.iana.org/assignments/ianaiftype-mib)
   -p (--port)       SNMP port (default 161)
   -I (--ifmib)      Agent supports IFMIB ifXTable. Do not use if
                     you don't know what this is.
   -n (--name)       the value should match the returned ifName
                     (Implies the use of -I)
   -w (--warn =i|w|c) ignore|warn|crit if the interface is dormant (default critical)
   -D (--admin-down =i|w|c) same for administratively down interfaces (default warning)
   -M (--maxmsgsize) Max message size - usefull only for v1 or v2c
   -t (--timeout)    seconds before the plugin times out (default=15)
   -V (--version)    Plugin version
   -h (--help)       usage help

 -k or -d or -T must be specified

Note: either -k or -d or -T must be specified and -d and -T are much more network
intensive.  Use it sparingly or not at all.  -n is used to match against
a much more descriptive ifName value in the IfXTable to verify that the
snmpkey has not changed to some other network interface after a reboot.`,
	"check_ifstatus": `check_ifstatus v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.

usage:
check_ifstatus -C <READCOMMUNITY> -p <PORT> -H <HOSTNAME>
Copyright (C) 2000 Christoph Kron
Updates 5/2002 Subhendu Ghosh
Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).


check_ifstatus plugin for Nagios monitors operational
status of each network interface on the target host

Usage:
   -H (--hostname)   Hostname to query - (required)
   -C (--community)  SNMP read community (defaults to public,
                     used with SNMP v1 and v2c
   -v (--snmp_version)  1 for SNMP v1 (default)
                        2 for SNMP v2c
                          SNMP v2c will use get_bulk for less overhead
                        3 for SNMPv3 (requires -U option)   -p (--port)       SNMP port (default 161)
   -I (--ifmib)      Agent supports IFMIB ifXTable.  For Cisco - this will provide
                     the descriptive name.  Do not use if you don't know what this is.
   -x (--exclude)    A comma separated list of ifType values that should be excluded
                     from the report (default for an empty list is PPP(23).
   -u (--unused_ports) A comma separated list of ifIndex values that should be excluded
                     from the report (default is an empty exclusion list).
                     See the IANAifType-MIB for a list of interface types.
   -L (--seclevel)   choice of "noAuthNoPriv", "authNoPriv", or	"authPriv"
   -U (--secname)    username for SNMPv3 context
   -c (--context)    SNMPv3 context name (default is empty string)
   -A (--authpass)   authentication password (cleartext ascii or localized key
                     in hex with 0x prefix generated by using "snmpkey" utility
                     auth password and authEngineID
   -a (--authproto)  Authentication protocol (MD5 or SHA1)
   -X (--privpass)   privacy password (cleartext ascii or localized key
                     in hex with 0x prefix generated by using "snmpkey" utility
                     privacy password and authEngineID
   -P (--privproto)  privacy protocol (DES or AES; default: DES)
   -M (--maxmsgsize) Max message size - usefull only for v1 or v2c
   -t (--timeout)    seconds before the plugin times out (default=15)
   -V (--version)    Plugin version
   -h (--help)       usage help

check_ifstatus v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.`,
	"check_imap": `check_imap v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests IMAP connections with the specified host (or unix socket).

Usage:
check_imap -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ircd": `check_ircd v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2000 Richard Mayhew/Karl DeBisschop

Perl Check IRCD plugin for Nagios

Usage: check_ircd -H <host> [-w <warn>] [-c <crit>] [-p <port>]

-H, --hostname=HOST
   Name or IP address of host to check
-w, --warning=INTEGER
   Number of connected users which generates a warning state (Default: 50)
-c, --critical=INTEGER
   Number of connected users which generates a critical state (Default: 100)
-p, --port=INTEGER
   Port that the ircd daemon is running on <host> (Default: 6667)
-v, --verbose
   Print extra debugging information`,
	"check_jabber": `check_jabber v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests JABBER connections with the specified host (or unix socket).

Usage:
check_jabber -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ldap": `check_ldap v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Didi Rieder (adrieder@sbox.tu-graz.ac.at)
Copyright (c) 2000-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>



Usage:
 check_ldap -H <host> -b <base_dn> [-p <port>] [-a <attr>] [-D <binddn>]
       [-P <password>] [-w <warn_time>] [-c <crit_time>] [-t timeout]
       [-2|-3] [-4|-6]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 389)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -a [--attr]
    ldap attribute to search (default: "(objectclass=*)"
 -b [--base]
    ldap base (eg. ou=my unit, o=my org, c=at
 -D [--bind]
    ldap bind DN (if required)
 -P [--pass]
    ldap password (if required)
 -T [--starttls]
    use starttls mechanism introduced in protocol version 3
 -S [--ssl]
    use ldaps (ldap v2 ssl method). this also sets the default port to 636
 -2 [--ver2]
    use ldap protocol version 2
 -3 [--ver3]
    use ldap protocol version 3
    (default protocol version: 2)
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Notes:
 If this plugin is called via 'check_ldaps', method 'STARTTLS' will be
 implied (using default port 389) unless --port=636 is specified. In that case
 'SSL on connect' will be used no matter how the plugin was called.
 This detection is deprecated, please use 'check_ldap' with the '--starttls' or '--ssl' flags
 to define the behaviour explicitly instead.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ldaps": `check_ldaps v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Didi Rieder (adrieder@sbox.tu-graz.ac.at)
Copyright (c) 2000-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>



Usage:
 check_ldaps -H <host> -b <base_dn> [-p <port>] [-a <attr>] [-D <binddn>]
       [-P <password>] [-w <warn_time>] [-c <crit_time>] [-t timeout]
       [-2|-3] [-4|-6]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 389)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -a [--attr]
    ldap attribute to search (default: "(objectclass=*)"
 -b [--base]
    ldap base (eg. ou=my unit, o=my org, c=at
 -D [--bind]
    ldap bind DN (if required)
 -P [--pass]
    ldap password (if required)
 -T [--starttls]
    use starttls mechanism introduced in protocol version 3
 -S [--ssl]
    use ldaps (ldap v2 ssl method). this also sets the default port to 636
 -2 [--ver2]
    use ldap protocol version 2
 -3 [--ver3]
    use ldap protocol version 3
    (default protocol version: 2)
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Notes:
 If this plugin is called via 'check_ldaps', method 'STARTTLS' will be
 implied (using default port 389) unless --port=636 is specified. In that case
 'SSL on connect' will be used no matter how the plugin was called.
 This detection is deprecated, please use 'check_ldap' with the '--starttls' or '--ssl' flags
 to define the behaviour explicitly instead.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_load": `check_load v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Felipe Gustavo de Almeida <galmeida@linux.ime.usp.br>
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests the current system load average.

Usage:
check_load [-r] -w WLOAD1,WLOAD5,WLOAD15 -c CLOAD1,CLOAD5,CLOAD15

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -w, --warning=WLOAD1,WLOAD5,WLOAD15
    Exit with WARNING status if load average exceeds WLOADn
 -c, --critical=CLOAD1,CLOAD5,CLOAD15
    Exit with CRITICAL status if load average exceed CLOADn
    the load average format is the same used by "uptime" and "w"
 -r, --percpu
    Divide the load averages by the number of CPUs (when possible)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_log": `check_log v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.

Usage: check_log -F logfile -O oldlog -q query
Usage: check_log --help
Usage: check_log --version

Log file pattern detector plugin for Nagios

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_mailq": `check_mailq v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2002 Subhendu Ghosh/Carlos Canau/Benjamin Schmid

Usage: check_mailq -w <warn> -c <crit> [-W <warn>] [-C <crit>] [-M <MTA>] [-t <timeout>] [-v verbose]

   Checks the number of messages in the mail queue (supports multiple sendmail queues, qmail)
   Feedback/patches to support non-sendmail mailqueue welcome

-w (--warning)   = Min. number of messages in queue to generate warning
-c (--critical)  = Min. number of messages in queue to generate critical alert ( w < c )
-W (--Warning)   = Min. number of messages for same domain in queue to generate warning
-C (--Critical)  = Min. number of messages for same domain in queue to generate critical alert ( W < C )
-t (--timeout)   = Plugin timeout in seconds (default = 15)
-M (--mailserver) = [ sendmail | qmail | postfix | exim ] (default = sendmail)
-h (--help)
-V (--version)
-v (--verbose)   = debugging output


Note: -w and -c are required arguments.  -W and -C are optional.
 -W and -C are applied to domains listed on the queues - both FROM and TO. (sendmail)
 -W and -C are applied message not yet preproccessed. (qmail)
 This plugin uses the system mailq command (sendmail) or qmail-stat (qmail)
 to look at the queues. Mailq can usually only be accessed by root or
 a TrustedUser. You will have to set appropriate permissions for the plugin to work.


Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_mrtg": `check_mrtg v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin will check either the average or maximum value of one of the
two variables recorded in an MRTG log file.


Usage:
check_mrtg -F log_file -a <AVG | MAX> -v variable -w warning -c critical
[-l label] [-u units] [-e expire_minutes] [-t timeout] [-v]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -F, --logfile=FILE
   The MRTG log file containing the data you want to monitor
 -e, --expires=MINUTES
   Minutes before MRTG data is considered to be too old
 -a, --aggregation=AVG|MAX
   Should we check average or maximum values?
 -v, --variable=INTEGER
   Which variable set should we inspect? (1 or 2)
 -w, --warning=INTEGER
   Threshold value for data to result in WARNING status
 -c, --critical=INTEGER
   Threshold value for data to result in CRITICAL status
 -l, --label=STRING
   Type label for data (Examples: Conns, "Processor Load", In, Out)
 -u, --units=STRING
   Option units label for data (Example: Packets/Sec, Errors/Sec,
   "Bytes Per Second", "%% Utilization")

 If the value exceeds the <vwl> threshold, a WARNING status is returned. If
 the value exceeds the <vcl> threshold, a CRITICAL status is returned.  If
 the data in the log file is older than <expire_minutes> old, a WARNING
 status is returned and a warning message is printed.

 This plugin is useful for monitoring MRTG data that does not correspond to
 bandwidth usage.  (Use the check_mrtgtraf plugin for monitoring bandwidth).
 It can be used to monitor any kind of data that MRTG is monitoring - errors,
 packets/sec, etc.  I use MRTG in conjuction with the Novell NLM that allows
 me to track processor utilization, user connections, drive space, etc and
 this plugin works well for monitoring that kind of data as well.

Notes:
 - This plugin only monitors one of the two variables stored in the MRTG log
   file.  If you want to monitor both values you will have to define two
   commands with different values for the <variable> argument.  Of course,
   you can always hack the code to make this plugin work for you...
 - MRTG stands for the Multi Router Traffic Grapher.  It can be downloaded from
   http://ee-staff.ethz.ch/~oetiker/webtools/mrtg/mrtg.html

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_mrtgtraf": `check_mrtgtraf v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin will check the incoming/outgoing transfer rates of a router,
switch, etc recorded in an MRTG log.  If the newest log entry is older
than <expire_minutes>, a WARNING status is returned. If either the
incoming or outgoing rates exceed the <icl> or <ocl> thresholds (in
Bytes/sec), a CRITICAL status results.  If either of the rates exceed
the <iwl> or <owl> thresholds (in Bytes/sec), a WARNING status results.


Usage check_mrtgtraf -F <log_file> -a <AVG | MAX> -w <warning_pair>
-c <critical_pair> [-e expire_minutes]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -F, --filename=STRING
    File to read log from
 -e, --expires=INTEGER
    Minutes after which log expires
 -a, --aggregation=(AVG|MAX)
    Test average or maximum
 -w, --warning
    Warning threshold pair <incoming>,<outgoing>
 -c, --critical
    Critical threshold pair <incoming>,<outgoing>

Notes:
 - MRTG stands for Multi Router Traffic Grapher. It can be downloaded from
   http://ee-staff.ethz.ch/~oetiker/webtools/mrtg/mrtg.html
 - While MRTG can monitor things other than traffic rates, this
   plugin probably won't work with much else without modification.
 - The calculated i/o rates are a little off from what MRTG actually
   reports.  I'm not sure why this is right now, but will look into it
   for future enhancements of this plugin.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_mysql": `check_mysql v1.5 (nagios-plugins 1.5)
Copyright (c) 1999-2011 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This program tests connections to a MySQL server


Usage:
 check_mysql [-d database] [-H host] [-P port] [-s socket]
       [-u user] [-p password] [-S] [-l] [-a cert] [-k key]
       [-C ca-cert] [-D ca-dir] [-L ciphers] [-f optfile] [-g group]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -P, --port=INTEGER
    Port number (default: 3306)
 -s, --socket=STRING
    Use the specified socket (has no effect if -H is used)
 -d, --database=STRING
    Check database with indicated name
 -f, --file=STRING
    Read from the specified client options file
 -g, --group=STRING
    Use a client options group
 -u, --username=STRING
    Connect using the indicated username
 -p, --password=STRING
    Use the indicated password to authenticate the connection
    ==> IMPORTANT: THIS FORM OF AUTHENTICATION IS NOT SECURE!!! <==
    Your clear-text password could be visible as a process table entry
 -S, --check-slave
    Check if the slave thread is running properly.
 -w, --warning
    Exit with WARNING status if slave server is more than INTEGER seconds
    behind master
 -c, --critical
    Exit with CRITICAL status if slave server is more then INTEGER seconds
    behind master
 -l, --ssl
    Use ssl encryptation
 -C, --ca-cert=STRING
    Path to CA signing the cert
 -a, --cert=STRING
    Path to SSL certificate
 -k, --key=STRING
    Path to private SSL key
 -D, --ca-dir=STRING
    Path to CA directory
 -L, --ciphers=STRING
    List of valid SSL ciphers

 There are no required arguments. By default, the local database is checked
 using the default unix socket. You can force TCP on localhost by using an
 IP address or FQDN ('localhost' will use the socket as well).

Notes:
 You must specify -p with an empty string to force an empty password,
 overriding any my.cnf settings.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_mysql_query": `check_mysql_query v1.5 (nagios-plugins 1.5)
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This program checks a query result against threshold levels


Usage:
 check_mysql_query -q SQL_query [-w warn] [-c crit] [-H host] [-P port] [-s socket]
       [-d database] [-u user] [-p password]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -q, --query=STRING
    SQL query to run. Only first column in first row will be read
 -w, --warning=RANGE
    Warning range (format: start:end). Alert if outside this range
 -c, --critical=RANGE
    Critical range
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -P, --port=INTEGER
    Port number (default: 3306)
 -s, --socket=STRING
    Use the specified socket (has no effect if -H is used)
 -d, --database=STRING
    Database to check
 -u, --username=STRING
    Username to login with
 -p, --password=STRING
    Password to login with
    ==> IMPORTANT: THIS FORM OF AUTHENTICATION IS NOT SECURE!!! <==
    Your clear-text password could be visible as a process table entry

 A query is required. The result from the query should be numeric.
 For extra security, create a user with minimal access.

Notes:
 You must specify -p with an empty string to force an empty password,
 overriding any my.cnf settings.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_nagios": `check_nagios v1.5 (nagios-plugins 1.5)
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks the status of the Nagios process on the local machine
The plugin will check to make sure the Nagios status log is no older than
the number of minutes specified by the expires option.
It also checks the process table for a process matching the command argument.


Usage:
check_nagios -F <status log file> -t <timeout_seconds> -e <expire_minutes> -C <process_string>

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -F, --filename=FILE
    Name of the log file to check
 -e, --expires=INTEGER
    Minutes aging after which logfile is considered stale
 -C, --command=STRING
    Substring to search for in process arguments
 -t, --timeout=INTEGER
    Timeout for the plugin in seconds
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Examples:
 check_nagios -t 20 -e 5 -F /usr/local/nagios/var/status.log -C /usr/local/nagios/bin/nagios

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_nntp": `check_nntp v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests NNTP connections with the specified host (or unix socket).

Usage:
check_nntp -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_nntps": `check_nntps v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests NNTPS connections with the specified host (or unix socket).

Usage:
check_nntps -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_nt": `check_nt v1.5 (nagios-plugins 1.5)
Copyright (c) 2000 Yves Rubin (rubiyz@yahoo.com)
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin collects data from the NSClient service running on a
Windows NT/2000/XP/2003 server.


Usage:
check_nt -H host -v variable [-p port] [-w warning] [-c critical]
[-l params] [-d SHOWALL] [-u] [-t timeout]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
Options:
 -H, --hostname=HOST
   Name of the host to check
 -p, --port=INTEGER
   Optional port number (default: 1248)
 -s, --secret=<password>
   Password needed for the request
 -w, --warning=INTEGER
   Threshold which will result in a warning status
 -c, --critical=INTEGER
   Threshold which will result in a critical status
 -t, --timeout=INTEGER
   Seconds before connection attempt times out (default:  -l, --params=<parameters>
   Parameters passed to specified check (see below) -d, --display={SHOWALL}
   Display options (currently only SHOWALL works) -u, --unknown-timeout
   Return UNKNOWN on timeouts10)
 -h, --help
   Print this help screen
 -V, --version
   Print version information
 -v, --variable=STRING
   Variable to check

Valid variables are:
 CLIENTVERSION = Get the NSClient version
  If -l <version> is specified, will return warning if versions differ.
 CPULOAD =
  Average CPU load on last x minutes.
  Request a -l parameter with the following syntax:
  -l <minutes range>,<warning threshold>,<critical threshold>.
  <minute range> should be less than 24*60.
  Thresholds are percentage and up to 10 requests can be done in one shot.
  ie: -l 60,90,95,120,90,95
 UPTIME =
  Get the uptime of the machine.
  No specific parameters. No warning or critical threshold
 USEDDISKSPACE =
  Size and percentage of disk use.
  Request a -l parameter containing the drive letter only.
  Warning and critical thresholds can be specified with -w and -c.
 MEMUSE =
  Memory use.
  Warning and critical thresholds can be specified with -w and -c.
 SERVICESTATE =
  Check the state of one or several services.
  Request a -l parameters with the following syntax:
  -l <service1>,<service2>,<service3>,...
  You can specify -d SHOWALL in case you want to see working services
  in the returned string.
 PROCSTATE =
  Check if one or several process are running.
  Same syntax as SERVICESTATE.
 COUNTER =
  Check any performance counter of Windows NT/2000.
	Request a -l parameters with the following syntax:
	-l "\\<performance object>\\counter","<description>
	The <description> parameter is optional and is given to a printf
  output command which requires a float parameter.
  If <description> does not include "%%", it is used as a label.
  Some examples:
  "Paging file usage is %%.2f %%%%"
  "%%.f %%%% paging file used."
 INSTANCES =
  Check any performance counter object of Windows NT/2000.
  Syntax: check_nt -H <hostname> -p <port> -v INSTANCES -l <counter object>
  <counter object> is a Windows Perfmon Counter object (eg. Process),
  if it is two words, it should be enclosed in quotes
  The returned results will be a comma-separated list of instances on
   the selected computer for that object.
  The purpose of this is to be run from command line to determine what instances
   are available for monitoring without having to log onto the Windows server
    to run Perfmon directly.
  It can also be used in scripts that automatically create Nagios service
   configuration files.
  Some examples:
  check_nt -H 192.168.1.1 -p 1248 -v INSTANCES -l Process

Notes:
 - The NSClient service should be running on the server to get any information
   (http://nsclient.ready2run.nl).
 - Critical thresholds should be lower than warning thresholds
 - Default port 1248 is sometimes in use by other services. The error
   output when this happens contains "Cannot map xxxxx to protocol number".
   One fix for this is to change the port to something else on check_nt
   and on the client service it's connecting to.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ntp": `check_ntp v1.5 (nagios-plugins 1.5)
Copyright (c) 2006 Sean Finney
Copyright (c) 2006-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks the selected ntp server


WARNING: check_ntp is deprecated. Please use check_ntp_peer or
check_ntp_time instead.

Usage:
 check_ntp -H <host> [-w <warn>] [-c <crit>] [-j <warn>] [-k <crit>] [-4|-6] [-v verbose]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 123)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -w, --warning=THRESHOLD
    Offset to result in warning status (seconds)
 -c, --critical=THRESHOLD
    Offset to result in critical status (seconds)
 -j, --jwarn=THRESHOLD
    Warning threshold for jitter
 -k, --jcrit=THRESHOLD
    Critical threshold for jitter
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Notes:
 See:
 http://nagiosplug.sourceforge.net/developer-guidelines.html#THRESHOLDFORMAT
 for THRESHOLD format and examples.

Examples:
 Normal offset check:
  ./check_ntp -H ntpserv -w 0.5 -c 1

 Check jitter too, avoiding critical notifications if jitter isn't available
 (See Notes above for more details on thresholds formats):
  ./check_ntp -H ntpserv -w 0.5 -c 1 -j -1:100 -k -1:200

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net

WARNING: check_ntp is deprecated. Please use check_ntp_peer or
check_ntp_time instead.`,
	"check_ntp_peer": `check_ntp_peer v1.5 (nagios-plugins 1.5)
Copyright (c) 2006 Sean Finney
Copyright (c) 2006-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks the selected ntp server


Usage:
 check_ntp_peer -H <host> [-4|-6] [-w <warn>] [-c <crit>] [-W <warn>] [-C <crit>]
       [-j <warn>] [-k <crit>] [-v verbose]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 123)
 -q, --quiet
    Returns UNKNOWN instead of CRITICAL or WARNING if server isn't synchronized
 -w, --warning=THRESHOLD
    Offset to result in warning status (seconds)
 -c, --critical=THRESHOLD
    Offset to result in critical status (seconds)
 -W, --swarn=THRESHOLD
    Warning threshold for stratum of server's synchronization peer
 -C, --scrit=THRESHOLD
    Critical threshold for stratum of server's synchronization peer
 -j, --jwarn=THRESHOLD
    Warning threshold for jitter
 -k, --jcrit=THRESHOLD
    Critical threshold for jitter
 -m, --twarn=THRESHOLD
    Warning threshold for number of usable time sources ("truechimers")
 -n, --tcrit=THRESHOLD
    Critical threshold for number of usable time sources ("truechimers")
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

This plugin checks an NTP server independent of any commandline
programs or external libraries.

Notes:
 Use this plugin to check the health of an NTP server. It supports
 checking the offset with the sync peer, the jitter and stratum. This
 plugin will not check the clock offset between the local host and NTP
 server; please use check_ntp_time for that purpose.

 See:
 http://nagiosplug.sourceforge.net/developer-guidelines.html#THRESHOLDFORMAT
 for THRESHOLD format and examples.

Examples:
 Simple NTP server check:
  ./check_ntp_peer -H ntpserv -w 0.5 -c 1

 Check jitter too, avoiding critical notifications if jitter isn't available
 (See Notes above for more details on thresholds formats):
  ./check_ntp_peer -H ntpserv -w 0.5 -c 1 -j -1:100 -k -1:200

 Only check the number of usable time sources ("truechimers"):
  ./check_ntp_peer -H ntpserv -m @5 -n @3

 Check only stratum:
  ./check_ntp_peer -H ntpserv -W 4 -C 6

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ntp_time": `check_ntp_time v1.5 (nagios-plugins 1.5)
Copyright (c) 2006 Sean Finney
Copyright (c) 2006-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks the clock offset with the ntp server


Usage:
 check_ntp_time -H <host> [-4|-6] [-w <warn>] [-c <crit>] [-v verbose]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 123)
 -q, --quiet
    Returns UNKNOWN instead of CRITICAL if offset cannot be found
 -w, --warning=THRESHOLD
    Offset to result in warning status (seconds)
 -c, --critical=THRESHOLD
    Offset to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

This plugin checks the clock offset between the local host and a
remote NTP server. It is independent of any commandline programs or
external libraries.

Notes:
 If you'd rather want to monitor an NTP server, please use
 check_ntp_peer.

 See:
 http://nagiosplug.sourceforge.net/developer-guidelines.html#THRESHOLDFORMAT
 for THRESHOLD format and examples.

Examples:
  ./check_ntp_time -H ntpserv -w 0.5 -c 1

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_nwstat": `check_nwstat v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin attempts to contact the MRTGEXT NLM running on a
Novell server to gather the requested system information.


Usage:
check_nwstat -H host [-p port] [-v variable] [-w warning] [-c critical] [-t timeout]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 9999)
 -v, --variable=STRING
   Variable to check.  Valid variables include:
    LOAD1     = 1 minute average CPU load
    LOAD5     = 5 minute average CPU load
    LOAD15    = 15 minute average CPU load
    CSPROCS   = number of current service processes (NW 5.x only)
    ABENDS    = number of abended threads (NW 5.x only)
    UPTIME    = server uptime
    LTCH      = percent long term cache hits
    CBUFF     = current number of cache buffers
    CDBUFF    = current number of dirty cache buffers
    DCB       = dirty cache buffers as a percentage of the total
    TCB       = dirty cache buffers as a percentage of the original
    OFILES    = number of open files
        VMF<vol>  = MB of free space on Volume <vol>
        VMU<vol>  = MB used space on Volume <vol>
        VMP<vol>  = MB of purgeable space on Volume <vol>
        VPF<vol>  = percent free space on volume <vol>
        VKF<vol>  = KB of free space on volume <vol>
        VPP<vol>  = percent purgeable space on volume <vol>
        VKP<vol>  = KB of purgeable space on volume <vol>
        VPNP<vol> = percent not yet purgeable space on volume <vol>
        VKNP<vol> = KB of not yet purgeable space on volume <vol>
        LRUM      = LRU sitting time in minutes
        LRUS      = LRU sitting time in seconds
        DSDB      = check to see if DS Database is open
        DSVER     = NDS version
        UPRB      = used packet receive buffers
        PUPRB     = percent (of max) used packet receive buffers
        SAPENTRIES = number of entries in the SAP table
        SAPENTRIES<n> = number of entries in the SAP table for SAP type <n>
        TSYNC     = timesync status
        LOGINS    = check to see if logins are enabled
        CONNS     = number of currently licensed connections
        NRMH	= NRM Summary Status
        NRMP<stat> = Returns the current value for a NRM health item
        NRMM<stat> = Returns the current memory stats from NRM
        NRMS<stat> = Returns the current Swapfile stats from NRM
        NSS1<stat> = Statistics from _Admin:Manage_NSS\GeneralStats.xml
        NSS3<stat> = Statistics from _Admin:Manage_NSS\NameCache.xml
        NSS4<stat> = Statistics from _Admin:Manage_NSS\FileStats.xml
        NSS5<stat> = Statistics from _Admin:Manage_NSS\ObjectCache.xml
        NSS6<stat> = Statistics from _Admin:Manage_NSS\Thread.xml
        NSS7<stat> = Statistics from _Admin:Manage_NSS\AuthorizationCache.xml
        NLM:<nlm> = check if NLM is loaded and report version
                    (e.g. NLM:TSANDS.NLM)

 -w, --warning=INTEGER
    Threshold which will result in a warning status
 -c, --critical=INTEGER
    Threshold which will result in a critical status
 -o, --osversion
    Include server version string in results
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)

Notes:
 - This plugin requres that the MRTGEXT.NLM file from James Drews' MRTG
   extension for NetWare be loaded on the Novell servers you wish to check.
   (available from http://www.engr.wisc.edu/~drews/mrtg/)
 - Values for critical thresholds should be lower than warning thresholds
   when the following variables are checked: VPF, VKF, LTCH, CBUFF, DCB,
   TCB, LRUS and LRUM.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_oracle": `check_oracle v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.

Usage:
  check_oracle --tns <Oracle Sid or Hostname/IP address>
  check_oracle --db <ORACLE_SID>
  check_oracle --login <ORACLE_SID>
  check_oracle --connect <ORACLE_SID>
  check_oracle --cache <ORACLE_SID> <USER> <PASS> <CRITICAL> <WARNING>
  check_oracle --tablespace <ORACLE_SID> <USER> <PASS> <TABLESPACE> <CRITICAL> <WARNING>
  check_oracle --oranames <Hostname>
  check_oracle --help
  check_oracle --version

Check Oracle status

--tns SID/IP Address
   Check remote TNS server
--db SID
   Check local database (search /bin/ps for PMON process) and check
   filesystem for sgadefORACLE_SID.dbf
--login SID
   Attempt a dummy login and alert if not ORA-01017: invalid username/password
--connect SID
   Attempt a login and alert if an ORA- error is returned
--cache
   Check local database for library and buffer cache hit ratios
       --->  Requires Oracle user/password and SID specified.
       		--->  Requires select on v_ and v_
--tablespace
   Check local database for tablespace capacity in ORACLE_SID
       --->  Requires Oracle user/password specified.
       		--->  Requires select on dba_data_files and dba_free_space
--oranames Hostname
   Check remote Oracle Names server
--help
   Print this help screen
--version
   Print version and license information

If the plugin doesn't work, check that the ORACLE_HOME environment
variable is set, that ORACLE_HOME/bin is in your PATH, and the
tnsnames.ora file is locatable and is properly configured.

When checking local database status your ORACLE_SID is case sensitive.

If you want to use a default Oracle home, add in your oratab file:
*:/opt/app/oracle/product/7.3.4:N

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_overcr": `check_overcr v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin attempts to contact the Over-CR collector daemon running on the
remote UNIX server in order to gather the requested system information.


Usage:
check_overcr -H host [-p port] [-v variable] [-w warning] [-c critical] [-t timeout]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 2000)
 -w, --warning=INTEGER
    Threshold which will result in a warning status
 -c, --critical=INTEGER
    Threshold which will result in a critical status
 -v, --variable=STRING
    Variable to check.  Valid variables include:
    LOAD1         = 1 minute average CPU load
    LOAD5         = 5 minute average CPU load
    LOAD15        = 15 minute average CPU load
    DPU<filesys>  = percent used disk space on filesystem <filesys>
    PROC<process> = number of running processes with name <process>
    NET<port>     = number of active connections on TCP port <port>
    UPTIME        = system uptime in seconds
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

This plugin requires that Eric Molitors' Over-CR collector daemon be
running on the remote server.
Over-CR can be downloaded from http://www.molitor.org/overcr
This plugin was tested with version 0.99.53 of the Over-CR collector

Notes:
 For the available options, the critical threshold value should always be
 higher than the warning threshold value, EXCEPT with the uptime variable

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_pgsql": `check_pgsql v1.5 (nagios-plugins 1.5)
Copyright (c) 1999-2011 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

Test whether a PostgreSQL Database is accepting connections.

Usage:
check_pgsql [-H <host>] [-P <port>] [-c <critical time>] [-w <warning time>]
 [-t <timeout>] [-d <database>] [-l <logname>] [-p <password>]
[-q <query>] [-C <critical query range>] [-W <warning query range>]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -P, --port=INTEGER
    Port number (default: 5432)
 -d, --database=STRING
    Database to check (default: template1) -l, --logname = STRING
    Login name of user
 -p, --password = STRING
    Password (BIG SECURITY ISSUE)
 -o, --option = STRING
    Connection parameters (keyword = value), see below
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -q, --query=STRING
    SQL query to run. Only first column in first row will be read
 -W, --query-warning=RANGE
    SQL query value to result in warning status (double)
 -C, --query-critical=RANGE
    SQL query value to result in critical status (double)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

 All parameters are optional.
 This plugin tests a PostgreSQL DBMS to determine whether it is active and
 accepting queries. In its current operation, it simply connects to the
 specified database, and then disconnects. If no database is specified, it
 connects to the template1 database, which is present in every functioning
 PostgreSQL DBMS.

 If a query is specified using the -q option, it will be executed after
 connecting to the server. The result from the query has to be numeric.
 Multiple SQL commands, separated by semicolon, are allowed but the result
 of the last command is taken into account only. The value of the first
 column in the first row is used as the check result.

 See the chapter "Monitoring Database Activity" of the PostgreSQL manual
 for details about how to access internal statistics of the database server.

 For a list of available connection parameters which may be used with the -o
 command line option, see the documentation for PQconnectdb() in the chapter
 "libpq - C Library" of the PostgreSQL manual. For example, this may be
 used to specify a service name in pg_service.conf to be used for additional
 connection parameters: -o 'service=<name>' or to specify the SSL mode:
 -o 'sslmode=require'.

 The plugin will connect to a local postmaster if no host is specified. To
 connect to a remote host, be sure that the remote postmaster accepts TCP/IP
 connections (start the postmaster with the -i option).

 Typically, the nagios user (unless the --logname option is used) should be
 able to connect to the database without a password. The plugin can also send
 a password, but no effort is made to obsure or encrypt the password.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ping": `check_ping v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

Use ping to check connection statistics for a remote host.

Usage:
check_ping -H <host_address> -w <wrta>,<wpl>% -c <crta>,<cpl>%
 [-p packets] [-t timeout] [-4|-6]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -H, --hostname=HOST
    host to ping
 -w, --warning=THRESHOLD
    warning threshold pair
 -c, --critical=THRESHOLD
    critical threshold pair
 -p, --packets=INTEGER
    number of ICMP ECHO packets to send (Default: 5)
 -L, --link
    show HTML in the plugin output (obsoleted by urlize)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)

THRESHOLD is <rta>,<pl>% where <rta> is the round trip average travel
time (ms) which triggers a WARNING or CRITICAL state, and <pl> is the
percentage of packet loss to trigger an alarm state.

This plugin uses the ping command to probe the specified host for packet loss
(percentage) and round trip average (milliseconds). It can produce HTML output
linking to a traceroute CGI contributed by Ian Cass. The CGI can be found in
the contrib area of the downloads section at http://www.nagios.org/

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_pop": `check_pop v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests POP connections with the specified host (or unix socket).

Usage:
check_pop -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_procs": `check_procs v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 2000-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

Checks all processes and generates WARNING or CRITICAL states if the specified
metric is outside the required threshold ranges. The metric defaults to number
of processes.  Search filters can be applied to limit the processes to check.


Usage:
check_procs -w <range> -c <range> [-m metric] [-s state] [-p ppid]
 [-u user] [-r rss] [-z vsz] [-P %cpu] [-a argument-array]
 [-C command] [-k] [-t timeout] [-v]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -w, --warning=RANGE
   Generate warning state if metric is outside this range
 -c, --critical=RANGE
   Generate critical state if metric is outside this range
 -m, --metric=TYPE
  Check thresholds against metric. Valid types:
  PROCS   - number of processes (default)
  VSZ     - virtual memory size
  RSS     - resident set memory size
  CPU     - percentage CPU
  ELAPSED - time elapsed in seconds
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Extra information. Up to 3 verbosity levels
 -T, --traditional
   Filter own process the traditional way by PID instead of /proc/pid/exe

Filters:
 -s, --state=STATUSFLAGS
   Only scan for processes that have, in the output of 'ps', one or
   more of the status flags you specify (for example R, Z, S, RS,
   RSZDT, plus others based on the output of your 'ps' command).
 -p, --ppid=PPID
   Only scan for children of the parent process ID indicated.
 -z, --vsz=VSZ
   Only scan for processes with VSZ higher than indicated.
 -r, --rss=RSS
   Only scan for processes with RSS higher than indicated.
 -P, --pcpu=PCPU
   Only scan for processes with PCPU higher than indicated.
 -u, --user=USER
   Only scan for processes with user name or ID indicated.
 -a, --argument-array=STRING
   Only scan for processes with args that contain STRING.
 --ereg-argument-array=STRING
   Only scan for processes with args that contain the regex STRING.
 -C, --command=COMMAND
   Only scan for exact matches of COMMAND (without path).
 -k, --no-kthreads
   Only scan for non kernel threads (works on Linux only).

RANGEs are specified 'min:max' or 'min:' or ':max' (or 'max'). If
specified 'max:min', a warning status will be generated if the
count is inside the specified range

This plugin checks the number of currently running processes and
generates WARNING or CRITICAL states if the process count is outside
the specified threshold ranges. The process count can be filtered by
process owner, parent process PID, current state (e.g., 'Z'), or may
be the total number of running processes

Examples:
 check_procs -w 2:2 -c 2:1024 -C portsentry
  Warning if not two processes with command name portsentry.
  Critical if < 2 or > 1024 processes

 check_procs -w 10 -a '/usr/local/bin/perl' -u root
  Warning alert if > 10 processes with command arguments containing
  '/usr/local/bin/perl' and owned by root

 check_procs -w 50000 -c 100000 --metric=VSZ
  Alert if VSZ of any processes over 50K or 100K

 check_procs -w 10 -c 20 --metric=CPU
  Alert if CPU of any processes over 10%% or 20%%

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_real": `check_real v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Pedro Leite <leite@cic.ua.pt>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests the REAL service on the specified host.


Usage:
check_real -H host [-e expect] [-p port] [-w warn] [-c crit] [-t timeout] [-v]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 554)
 -u, --url=STRING
    Connect to this url
 -e, --expect=STRING
String to expect in first line of server response (default: RTSP/1.)
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

This plugin will attempt to open an RTSP connection with the host.
Successul connects return STATE_OK, refusals and timeouts return
STATE_CRITICAL, other errors return STATE_UNKNOWN.  Successful connects,
but incorrect reponse messages from the host result in STATE_WARNING return
values.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_rpc": `check_rpc v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2002 Karl DeBisschop/Truongchinh Nguyen/Subhendu Ghosh

Check if a rpc service is registered and running using
      rpcinfo -H host -C rpc_command

Usage:
 check_rpc -H host -C rpc_command [-p port] [-c program_version] [-u|-t] [-v]
 check_rpc [-h | --help]
 check_rpc [-V | --version]

  <host>          The server providing the rpc service
  <rpc_command>   The program name (or number).
  <program_version> The version you want to check for (one or more)
                    Should prevent checks of unknown versions being syslogged
                    e.g. 2,3,6 to check v2, v3, and v6
  [-u | -t]       Test UDP or TCP
  [-v]            Verbose
  [-v -v]         Verbose - will print supported programs and numbers

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_rta_multi": `Copyright (c) 2005 Andreas Ericsson <ae@op5.se>
Copyright (c) 2005-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>



Usage:
 check_rta_multi [options] [-H] host1 host2 hostN

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H
    specify a target
 -w
    warning threshold (currently 200,000ms,40%)
 -c
    critical threshold (currently 500,000ms,80%)
 -s
    specify a source IP address or device name
 -n
    number of packets to send (currently 5)
 -i
    max packet interval (currently 50,000ms)
 -I
    max target interval (currently 0,000ms)
 -m
    number of alive hosts required for success
 -l
    TTL on outgoing packets (currently 0)
 -t
    timeout value (seconds, currently  10)
 -b
    Number of icmp data bytes to send
    Packet size will be data bytes + icmp header (currently 68 + 8)
 -v
    verbose

Notes:
 The -H switch is optional. Naming a host (or several) to check is not.

 Threshold format for -w and -c is 200.25,60% for 200.25 msec RTA and 60%
 packet loss.  The default values should work well for most users.
 You can specify different RTA factors using the standardized abbreviations
 us (microseconds), ms (milliseconds, default) or just plain s for seconds.

 The -v switch can be specified several times for increased verbosity.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_sensors": `check_sensors v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.

Usage: check_sensors [--ignore-fault]

This plugin checks hardware status using the lm_sensors package.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
	"check_simap": `check_simap v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests SIMAP connections with the specified host (or unix socket).

Usage:
check_simap -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_smtp": `check_smtp v1.5 (nagios-plugins 1.5)
Copyright (c) 1999-2001 Ethan Galstad <nagios@nagios.org>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin will attempt to open an SMTP connection with the host.


Usage:
check_smtp -H host [-p port] [-4|-6] [-e expect] [-C command] [-R response] [-f from addr]
[-A authtype -U authuser -P authpass] [-w warn] [-c crit] [-t timeout] [-q]
[-F fqdn] [-S] [-D warn days cert expire[,crit days cert expire]] [-v]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 25)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -e, --expect=STRING
    String to expect in first line of server response (default: '220')
 -C, --command=STRING
    SMTP command (may be used repeatedly)
 -R, --response=STRING
    Expected response to command (may be used repeatedly)
 -f, --from=STRING
    FROM-address to include in MAIL command, required by Exchange 2000
 -F, --fqdn=STRING
    FQDN used for HELO
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
 -S, --starttls
    Use STARTTLS for the connection.
 -A, --authtype=STRING
    SMTP AUTH type to check (default none, only LOGIN supported)
 -U, --authuser=STRING
    SMTP AUTH username
 -P, --authpass=STRING
    SMTP AUTH password
 -q, --ignore-quit-failure
    Ignore failure when sending QUIT command to server
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Successul connects return STATE_OK, refusals and timeouts return
STATE_CRITICAL, other errors return STATE_UNKNOWN.  Successful
connects, but incorrect reponse messages from the host result in
STATE_WARNING return values.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_snmp": `check_snmp v1.5 (nagios-plugins 1.5)
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

Check status of remote machines and obtain system information via SNMP


Usage:
check_snmp -H <ip_address> -o <OID> [-w warn_range] [-c crit_range]
[-C community] [-s string] [-r regex] [-R regexi] [-t timeout] [-e retries]
[-l label] [-u units] [-p port-number] [-d delimiter] [-D output-delimiter]
[-m miblist] [-P snmp version] [-L seclevel] [-U secname] [-a authproto]
[-A authpasswd] [-x privproto] [-X privpasswd]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 161)
 -n, --next
    Use SNMP GETNEXT instead of SNMP GET
 -P, --protocol=[1|2c|3]
    SNMP protocol version
 -L, --seclevel=[noAuthNoPriv|authNoPriv|authPriv]
    SNMPv3 securityLevel
 -a, --authproto=[MD5|SHA]
    SNMPv3 auth proto
 -x, --privproto=[DES|AES]
    SNMPv3 priv proto (default DES)
 -C, --community=STRING
    Optional community string for SNMP communication (default is "public")
 -U, --secname=USERNAME
    SNMPv3 username
 -A, --authpassword=PASSWORD
    SNMPv3 authentication password
 -X, --privpasswd=PASSWORD
    SNMPv3 privacy password
 -o, --oid=OID(s)
    Object identifier(s) or SNMP variables whose value you wish to query
 -m, --miblist=STRING
    List of MIBS to be loaded (default = none if using numeric OIDs or 'ALL'
    for symbolic OIDs.)
 -d, --delimiter=STRING
    Delimiter to use when parsing returned data. Default is "="
    Any data on the right hand side of the delimiter is considered
    to be the data that should be used in the evaluation.
 -w, --warning=THRESHOLD(s)
    Warning threshold range(s)
 -c, --critical=THRESHOLD(s)
    Critical threshold range(s)
 --rate
    Enable rate calculation. See 'Rate Calculation' below
 --rate-multiplier
    Converts rate per second. For example, set to 60 to convert to per minute
 --offset=OFFSET
    Add/substract the specified OFFSET to numeric sensor data
 -s, --string=STRING
    Return OK state (for that OID) if STRING is an exact match
 -r, --ereg=REGEX
    Return OK state (for that OID) if extended regular expression REGEX matches
 -R, --eregi=REGEX
    Return OK state (for that OID) if case-insensitive extended REGEX matches
 --invert-search
    Invert search result (CRITICAL if found)
 -l, --label=STRING
    Prefix label for output from plugin
 -u, --units=STRING
    Units label(s) for output data (e.g., 'sec.').
 -D, --output-delimiter=STRING
    Separates output on multiple OID requests
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -e, --retries=INTEGER
    Number of retries to be used in the requests
 -O, --perf-oids
    Label performance data with OIDs instead of --label's
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

This plugin uses the 'snmpget' command included with the NET-SNMP package.
if you don't have the package installed, you will need to download it from
http://net-snmp.sourceforge.net before you can use this plugin.

Notes:
 - Multiple OIDs (and labels) may be indicated by a comma or space-delimited
   list (lists with internal spaces must be quoted).
 - See:
 http://nagiosplug.sourceforge.net/developer-guidelines.html#THRESHOLDFORMAT
 for THRESHOLD format and examples.
 - When checking multiple OIDs, separate ranges by commas like '-w 1:10,1:,:20'
 - Note that only one string and one regex may be checked at present
 - All evaluation methods other than PR, STR, and SUBSTR expect that the value
   returned from the SNMP query is an unsigned integer.

Rate Calculation:
 In many places, SNMP returns counters that are only meaningful when
 calculating the counter difference since the last check. check_snmp
 saves the last state information in a file so that the rate per second
 can be calculated. Use the --rate option to save state information.
 On the first run, there will be no prior state - this will return with OK.
 The state is uniquely determined by the arguments to the plugin, so
 changing the arguments will create a new state file.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_spop": `check_spop v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests SPOP connections with the specified host (or unix socket).

Usage:
check_spop -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ssh": `check_ssh v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Remi Paulmier <remi@sinfomic.fr>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

Try to connect to an SSH server at specified server and port


Usage:
check_ssh  [-4|-6] [-t <timeout>] [-r <remote version>] [-p <port>] <host>

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 22)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -r, --remote-version=STRING
    Warn if string doesn't match expected server version (ex: OpenSSH_3.9p1)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ssmtp": `check_ssmtp v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests SSMTP connections with the specified host (or unix socket).

Usage:
check_ssmtp -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_swap": `check_swap v1.5 (nagios-plugins 1.5)
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

Check swap space on local machine.


Usage:
check_swap [-av] -w <percent_free>% -c <percent_free>%
check_swap [-av] -w <bytes_free> -c <bytes_free>

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -w, --warning=INTEGER
    Exit with WARNING status if less than INTEGER bytes of swap space are free
 -w, --warning=PERCENT%%
    Exit with WARNING status if less than PERCENT of swap space is free
 -c, --critical=INTEGER
    Exit with CRITICAL status if less than INTEGER bytes of swap space are free
 -c, --critical=PERCENT%%
    Exit with CRITCAL status if less than PERCENT of swap space is free
 -a, --allswaps
    Conduct comparisons for all swap partitions, one by one
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Notes:
 On AIX, if -a is specified, uses lsps -a, otherwise uses lsps -s.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_tcp": `check_tcp v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests TCP connections with the specified host (or unix socket).

Usage:
check_tcp -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_time": `check_time v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad
Copyright (c) 1999-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin will check the time on the specified host.


Usage:
check_time -H <host_address> [-p port] [-u] [-w variance] [-c variance]
 [-W connect_time] [-C connect_time] [-t timeout]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 37)
 -u, --udp
   Use UDP to connect, not TCP
 -w, --warning-variance=INTEGER
   Time difference (sec.) necessary to result in a warning status
 -c, --critical-variance=INTEGER
   Time difference (sec.) necessary to result in a critical status
 -W, --warning-connect=INTEGER
   Response time (sec.) necessary to result in warning status
 -C, --critical-connect=INTEGER
   Response time (sec.) necessary to result in critical status
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_udp": `check_udp v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad <nagios@nagios.org>
Copyright (c) 1999-2008 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests UDP connections with the specified host (or unix socket).

Usage:
check_udp -H host -p port [-w <warning time>] [-c <critical time>] [-s <send string>]
[-e <expect string>] [-q <quit string>][-m <maximum bytes>] [-d <delay>]
[-t <timeout seconds>] [-r <refuse state>] [-M <mismatch state>] [-v] [-4|-6] [-j]
[-D <warn days cert expire>[,<crit days cert expire>]] [-S <use SSL>] [-E]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: none)
 -4, --use-ipv4
    Use IPv4 connection
 -6, --use-ipv6
    Use IPv6 connection
 -E, --escape
    Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option
    Default: nothing added to send, \r\n added to end of quit
 -s, --send=STRING
    String to send to the server
 -e, --expect=STRING
    String to expect in server response (may be repeated)
 -A, --all
    All expect strings need to occur in server response. Default is any
 -q, --quit=STRING
    String to send server to initiate a clean close of the connection
 -r, --refuse=ok|warn|crit
    Accept TCP refusals with states ok, warn, crit (default: crit)
 -M, --mismatch=ok|warn|crit
    Accept expected string mismatches with states ok, warn, crit (default: warn)
 -j, --jail
    Hide output from TCP socket
 -m, --maxbytes=INTEGER
    Close connection once more than this number of bytes are received
 -d, --delay=INTEGER
    Seconds to wait between sending string and polling for response
 -D, --certificate=INTEGER[,INTEGER]
    Minimum number of days a certificate has to be valid.
    1st is #days for warning, 2nd is critical (if not specified - 0).
 -S, --ssl
    Use SSL for the connection.
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)
 -v, --verbose
    Show details for command-line debugging (Nagios may truncate output)

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_ups": `check_ups v1.5 (nagios-plugins 1.5)
Copyright (c) 2000 Tom Shields
Copyright (c) 2004 Alain Richard <alain.richard@equation.fr>
Copyright (c) 2004 Arnaud Quette <arnaud.quette@mgeups.com>
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin tests the UPS service on the specified host. Network UPS Tools
from www.networkupstools.org must be running for this plugin to work.


Usage:
check_ups -H host -u ups [-p port] [-v variable] [-w warn_value] [-c crit_value] [-to to_sec] [-T]

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -H, --hostname=ADDRESS
    Host name, IP Address, or unix socket (must be an absolute path)
 -p, --port=INTEGER
    Port number (default: 3493)
 -u, --ups=STRING
    Name of UPS
 -T, --temperature
    Output of temperatures in Celsius
 -v, --variable=STRING
    Valid values for STRING are LINE, TEMP, BATTPCT or LOADPCT
 -w, --warning=DOUBLE
    Response time to result in warning status (seconds)
 -c, --critical=DOUBLE
    Response time to result in critical status (seconds)
 -t, --timeout=INTEGER
    Seconds before connection times out (default: 10)

This plugin attempts to determine the status of a UPS (Uninterruptible Power
Supply) on a local or remote host. If the UPS is online or calibrating, the
plugin will return an OK state. If the battery is on it will return a WARNING
state. If the UPS is off or has a low battery the plugin will return a CRITICAL
state.

Notes:
 You may also specify a variable to check (such as temperature, utility voltage,
 battery load, etc.) as well as warning and critical thresholds for the value
 of that variable.  If the remote host has multiple UPS that are being monitored
 you will have to use the --ups option to specify which UPS to check.

 This plugin requires that the UPSD daemon distributed with Russell Kroll's
 Network UPS Tools be installed on the remote host. If you do not have the
 package installed on your system, you can download it from
 http://www.networkupstools.org

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_users": `check_users v1.5 (nagios-plugins 1.5)
Copyright (c) 1999 Ethan Galstad
Copyright (c) 2000-2007 Nagios Plugin Development Team
	<nagiosplug-devel@lists.sourceforge.net>

This plugin checks the number of users currently logged in on the local
system and generates an error if the number exceeds the thresholds specified.


Usage:
check_users -w <users> -c <users>

Options:
 -h, --help
    Print detailed help screen
 -V, --version
    Print version information
 --extra-opts=[section][@file]
    Read options from an ini file. See http://nagiosplugins.org/extra-opts
    for usage and examples.
 -w, --warning=INTEGER
    Set WARNING status if more than INTEGER users are logged in
 -c, --critical=INTEGER
    Set CRITICAL status if more than INTEGER users are logged in

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net`,
	"check_wave": `check_wave v1.5 (nagios-plugins 1.5)
The nagios plugins come with ABSOLUTELY NO WARRANTY. You may redistribute
copies of the plugins under the terms of the GNU General Public License.
For more information about these matters, see the file named COPYING.
Copyright (c) 2000 Jeffery Blank/Karl DeBisschop

Usage: check_wave -H <host> [-w <warn>] [-c <crit>]

<warn> = Signal strength at which a warning message will be generated.
<crit> = Signal strength at which a critical message will be generated.

Send email to nagios-users@lists.sourceforge.net if you have questions
regarding use of this software. To submit patches or suggest improvements,
send email to nagiosplug-devel@lists.sourceforge.net.
Please include version information with all correspondence (when possible,
use output from the --version option of the plugin itself).`,
}
