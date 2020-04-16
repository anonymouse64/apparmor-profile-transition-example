# Example of transitioning AppArmor profile from Go program

This is an example of how to transition to a new AppArmor Profile in a Go program. 

This example assumes the specified AppArmor profile is already loaded in the system, that can be done with something like this:

```bash
$ cat > profile.txt << EOF
#include <tunables/global>

profile example-child-profile flags=(attach_disconnected,mediate_deleted) {
  # NOTE: this is not considered a "secure" profile, just an example
  #include <abstractions/base>

  # all network and file access
  network,
  file,

  # deny some privileged things in /proc and /sys
  deny @{PROC}/* w,   # deny write for all files directly in /proc (not in a subdir)
  # deny write to files not in /proc/<number>/** or /proc/sys/**
  deny @{PROC}/{[^1-9],[^1-9][^0-9],[^1-9s][^0-9y][^0-9s],[^1-9][^0-9][^0-9][^0-9]*}/** w,
  deny @{PROC}/sys/[^k]** w,  # deny /proc/sys except /proc/sys/k* (effectively /proc/sys/kernel)
  deny @{PROC}/sys/kernel/{?,??,[^s][^h][^m]**} w,  # deny everything except shm* in /proc/sys/kernel/
  deny @{PROC}/sysrq-trigger rwklx,
  deny @{PROC}/mem rwklx,
  deny @{PROC}/kmem rwklx,
  deny @{PROC}/kcore rwklx,
  deny /sys/[^f]*/** wklx,
  deny /sys/f[^s]*/** wklx,
  deny /sys/fs/[^c]*/** wklx,
  deny /sys/fs/c[^g]*/** wklx,
  deny /sys/fs/cg[^r]*/** wklx,
  deny /sys/firmware/efi/efivars/** rwklx,
  deny /sys/kernel/security/** rwklx,

  # also deny mounting things
  deny mount,
}

EOF
$ sudo apparmor_parser -r profile.txt
```

Then you would transition to it with this program like follows:

```bash
$ ./apparmor-profile-transition-example example-child-profile /bin/bash
$ cat /sys/kernel/security/lsm
cat: /sys/kernel/security/lsm: Permission denied
$ 
```