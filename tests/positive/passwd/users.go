// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package passwd

import (
	"github.com/coreos/ignition/tests/types"
)

func Add_passwd_users() types.Test {
	var mntDevices []types.MntDevice

	name := "Adding users"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	mntDevices = nil
	config := `{
		"ignition": {
			"version": "2.0.0"
		},
		"passwd": {
			"users": [{
					"name": "test",
					"create": {},
					"passwordHash": "zJW/EKqqIk44o",
					"sshAuthorizedKeys": [
						"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDBRZPFJNOvQRfokigTtl0IBi71LHZrFOk4EJ3Zowtk/bX5uIVai0Cd4+hqlocYL10idgtFBH28skeKfsmHwgS9XwOvP+g+kqAl7yCz8JEzIUzl1fxNZDToi0jA3B5MwXkpt+IWfnabwi2cRZhlzrz9rO+eExu5s3NfaRmmmCYrjCJIRPKSCrW8U0n9fVSbX4PDdMXVmH7r+t8MtR8523vCbakFR/Y0YIqkPVdfuUXHh9rDCdH4B7mt7nYX2LWQXGUvmI13mgQoy04ifkaR3ImuOMp3Y1J1gm6clO74IMCq/sK9+XJhbxMPPHUoUJ2EwbaG7Dbh3iqz47e9oVki4gIH stephenlowrie@localhost.localdomain"
					]
				},
				{
					"name": "jenkins",
					"create": {
						"uid": 1020
					}
				}
			]
		}
	}`
	in[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name: "passwd",
				Path: "etc",
			},
			Contents: []string{"root:x:0:0:root:/root:/bin/bash\ncore:x:500:500:CoreOS Admin:/home/core:/bin/bash\nsystemd-coredump:x:998:998:systemd Core Dumper:/:/sbin/nologin\nfleet:x:253:253::/:/sbin/nologin\n"},
		},
		{
			Node: types.Node{
				Name: "shadow",
				Path: "etc",
			},
			Contents: []string{"root:*:15887:0:::::\ncore:*:15887:0:::::\nsystemd-coredump:!!:17301::::::\nfleet:!!:17301::::::\n"},
		},
		{
			Node: types.Node{
				Name: "group",
				Path: "etc",
			},
			Contents: []string{"root:x:0:root\nwheel:x:10:root,core\nsudo:x:150:\ndocker:x:233:core\nsystemd-coredump:x:998:\nfleet:x:253:core\ncore:x:500:\nrkt-admin:x:999:\nrkt:x:251:core\n"},
		},
		{
			Node: types.Node{
				Name: "gshadow",
				Path: "etc",
			},
			Contents: []string{"root:*::root\nusers:*::\nsudo:*::\nwheel:*::root,core\nsudo:*::\ndocker:*::core\nsystemd-coredump:!!::\nfleet:!!::core\nrkt-admin:!!::\nrkt:!!::core\ncore:*::\n"},
		},
		{
			Node: types.Node{
				Name: "nsswitch.conf",
				Path: "etc",
			},
			Contents: []string{"# /etc/nsswitch.conf:\n\npasswd:      files\nshadow:      files\ngroup:       files\n\nhosts:       files dns myhostname\nnetworks:    files dns\n\nservices:    files\nprotocols:   files\nrpc:         files\n\nethers:      files\nnetmasks:    files\nnetgroup:    files\nbootparams:  files\nautomount:   files\naliases:     files\n"},
		},
		{
			Node: types.Node{
				Name: "login.defs",
				Path: "etc",
			},
			Contents: []string{`#
# Please note that the parameters in this configuration file control the
# behavior of the tools from the shadow-utils component. None of these
# tools uses the PAM mechanism, and the utilities that use PAM (such as the
# passwd command) should therefore be configured elsewhere. Refer to
# /etc/pam.d/system-auth for more information.
#

# *REQUIRED*
#   Directory where mailboxes reside, _or_ name of file, relative to the
#   home directory.  If you _do_ define both, MAIL_DIR takes precedence.
#   QMAIL_DIR is for Qmail
#
#QMAIL_DIR	Maildir
MAIL_DIR	/var/spool/mail
#MAIL_FILE	.mail

# Password aging controls:
#
#	PASS_MAX_DAYS	Maximum number of days a password may be used.
#	PASS_MIN_DAYS	Minimum number of days allowed between password changes.
#	PASS_MIN_LEN	Minimum acceptable password length.
#	PASS_WARN_AGE	Number of days warning given before a password expires.
#
PASS_MAX_DAYS	99999
PASS_MIN_DAYS	0
PASS_MIN_LEN	5
PASS_WARN_AGE	7

#
# Min/max values for automatic uid selection in useradd
#
UID_MIN                  1000
UID_MAX                 60000
# System accounts
SYS_UID_MIN               201
SYS_UID_MAX               999

#
# Min/max values for automatic gid selection in groupadd
#
GID_MIN                  1000
GID_MAX                 60000
# System accounts
SYS_GID_MIN               201
SYS_GID_MAX               999

#
# If defined, this command is run when removing a user.
# It should remove any at/cron/print jobs etc. owned by
# the user to be removed (passed as the first argument).
#
#USERDEL_CMD	/usr/sbin/userdel_local

#
# If useradd should create home directories for users by default
# On RH systems, we do. This option is overridden with the -m flag on
# useradd command line.
#
CREATE_HOME	yes

# The permission mask is initialized to this value. If not specified,
# the permission mask will be initialized to 022.
UMASK           077

# This enables userdel to remove user groups if no members exist.
#
USERGROUPS_ENAB yes

# Use SHA512 to encrypt password.
ENCRYPT_METHOD SHA512
`},
		},
	})
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name: "passwd",
				Path: "etc",
			},
			Contents: []string{"root:x:0:0:root:/root:/bin/bash\ncore:x:500:500:CoreOS Admin:/home/core:/bin/bash\nsystemd-coredump:x:998:998:systemd Core Dumper:/:/sbin/nologin\nfleet:x:253:253::/:/sbin/nologin\ntest:x:1000:1000::/home/test:/bin/bash\njenkins:x:1020:1001::/home/jenkins:/bin/bash\n"},
		},
		{
			Node: types.Node{
				Name: "group",
				Path: "etc",
			},
			Contents: []string{"root:x:0:root\nwheel:x:10:root,core\nsudo:x:150:\ndocker:x:233:core\nsystemd-coredump:x:998:\nfleet:x:253:core\ncore:x:500:\nrkt-admin:x:999:\nrkt:x:251:core\ntest:x:1000:\njenkins:x:1001:\n"},
		},
		{
			Node: types.Node{
				Name: "shadow",
				Path: "etc",
			},
			Contents: []string{"root:*:15887:0:::::\ncore:*:15887:0:::::\nsystemd-coredump:!!:17301::::::\nfleet:!!:17301::::::\ntest:zJW/EKqqIk44o:17331:0:99999:7:::\njenkins:*:17331:0:99999:7:::\n"},
		},
		{
			Node: types.Node{
				Name: "gshadow",
				Path: "etc",
			},
			Contents: []string{"root:*::root\nusers:*::\nsudo:*::\nwheel:*::root,core\nsudo:*::\ndocker:*::core\nsystemd-coredump:!!::\nfleet:!!::core\nrkt-admin:!!::\nrkt:!!::core\ncore:*::\ntest:!::\njenkins:!::\n"},
		},
		{
			Node: types.Node{
				Name:  "authorized_keys",
				Path:  "home/test/.ssh",
				User:  1000,
				Group: 1000,
			},
			Contents: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDBRZPFJNOvQRfokigTtl0IBi71LHZrFOk4EJ3Zowtk/bX5uIVai0Cd4+hqlocYL10idgtFBH28skeKfsmHwgS9XwOvP+g+kqAl7yCz8JEzIUzl1fxNZDToi0jA3B5MwXkpt+IWfnabwi2cRZhlzrz9rO+eExu5s3NfaRmmmCYrjCJIRPKSCrW8U0n9fVSbX4PDdMXVmH7r+t8MtR8523vCbakFR/Y0YIqkPVdfuUXHh9rDCdH4B7mt7nYX2LWQXGUvmI13mgQoy04ifkaR3ImuOMp3Y1J1gm6clO74IMCq/sK9+XJhbxMPPHUoUJ2EwbaG7Dbh3iqz47e9oVki4gIH stephenlowrie@localhost.localdomain\n\n"},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}
