pip list --outdated | awk 'NR>2 {print $1}' | xargs -n1 pip install --upgrade
