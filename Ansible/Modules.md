# Modules

### user

| **Parameter**              | **Choices/**Defaults    | **Comments**                                                 |
| -------------------------- | ----------------------- | ------------------------------------------------------------ |
| **name** string / required |                         | Name of the user to create, remove or modify. aliases: user  |
| **append** boolean         | **Choices:** **no** ← yes | If `yes`, add the user to the groups specified in `groups`.If `no`, user will only be added to the groups specified in `groups`, removing them from all other groups.Mutually exclusive with `local` |
| **state** string | **Choices:** absent **present** ← | Whether the account should exist or not, taking action if the state is different from what is stated. |
| **create_home** boolean | **Choices:** no **yes** ← | Unless set to `no`, a home directory will be made for the user when the account is created or if the home directory does not exist.Changed from `createhome` to `create_home` in Ansible 2.5. aliases: createhome |
| **shell** string |  | Optionally set the user's shell.On macOS, before Ansible 2.5, the default shell for non-system users was `/usr/bin/false`. Since Ansible 2.5, the default shell for non-system users on macOS is `/bin/bash`.On other operating systems, the default shell is determined by the underlying tool being used. See Notes for details. |

### lineinfile

| Parameter                                             | Comments                                                     |
| ----------------------------------------------------- | ------------------------------------------------------------ |
| **path** aliases: dest, destfile, namepath / required | The file to modify.Before Ansible 2.3 this option was only usable as *dest*, *destfile* and *name*. |
| **line** aliases: value string                        | The line to insert/replace into the file.Required for `state=present`.If `backrefs` is set, may contain backreferences that will get expanded with the `regexp` capture groups if the regexp matches. |
| **validate** string                                   | The validation command to run before copying the updated file into the final destination.A temporary file path is used to validate, passed in through ‘%s’ which must be present as in the examples below.Also, the command is passed securely so shell features such as expansion and pipes will not work.For an example on how to handle more complex validation than what this option provides, see [Complex configuration validation](https://docs.ansible.com/ansible/devel/reference_appendices/faq.html). |

### authorized_key

| Parameter                  | Comments                                                     |
| -------------------------- | ------------------------------------------------------------ |
| **user** string / required | The username on the remote host whose authorized_keys file will be modified. |
| **key** string / required  | The SSH public key(s), as a string or (since Ansible 1.9) url (https://github.com/username.keys). |

