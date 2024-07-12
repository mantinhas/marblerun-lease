show_help() {
	echo "Usage: maniedit [command] [object-version]"
	echo "                                            "
	echo "Commands:"
	echo "  save		save working version as [object-version]"
	echo "  set		set [object-version] as working version"
	echo "  help		get help message"
	echo "                                            "
	echo "Versions:"
	echo "  gramine		gramine manifest"
	echo "  marblerun	marblerun manifest"
}


if [ $# -ne 2 ]; then
	    show_help
	        exit 0
fi

case "$2" in
	gramine)
		saved_file=python.manifest.template.gramine
		;;
	marblerun)
		saved_file=python.manifest.template.marblerun
		;;
	*)
		show_help
		exit 0
		;;
esac

case "$1" in
	set)
		cp python.manifest.template python.manifest.template.autobak
		cp $saved_file python.manifest.template 
		;;
	save)
		cp $saved_file $saved_file.autobak
		cp python.manifest.template $saved_file
		;;
esac
