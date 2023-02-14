terraform init

# The original terraform graph is...well it includes a lot of data
terraform graph > 1.dot
cat 1.dot | dot -Tsvg > 1.svg

# I tried out this one, and found it was pretty ok!
# I wish the graphviz had color though, like it does when it renders html
# cat 1.dot | terraform-graph-beautifier -output-type graphviz > 2.dot
# cat 2.dot | dot -Tsvg > 2.svg

# Another take, but I don't think this had any effect
# cat 1.dot | terraform-graph-beautifier \
#     --exclude="module.root.provider" \
#     -output-type graphviz > 3.dot
# cat 3.dot | dot -Tsvg > 3.svg

# This one was...oof, didn't work out of the box
# docker run --rm -it -p 5000:5000 \
#   -v $(pwd):/data:ro \
#   --security-opt apparmor:unconfined \
#   --cap-add=SYS_ADMIN \
#   28mm/blast-radius

# I tried running it as a python module.
# This failed with a big ol error. Probably related to TF versions.
# pip3 install blast-radius
# blast-radius --serve .

# So I forked it and tried correcting the Dockerfile. Forced it into TF 0.14.5
# No dice, same error. Could maybe fix it, or use someone else's fork of it...
# docker run --rm -it -p 5000:5000 \
#   -v $(pwd):/data:ro \
#   --security-opt apparmor:unconfined \
#   --cap-add=SYS_ADMIN \
#   spilliams/blast-radius

# Still searching for that perfect grapher...
