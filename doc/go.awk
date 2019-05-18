BEGIN { show = 0 ; print "/*" }

/^\-/ { trim = 1 ; print "" }

/^Package/ { show = 1 }

!NF { trim = 0 }

trim { sub("^ +", "", $0) }

show { print $0 }

END { print "*/\npackage " package_name }
