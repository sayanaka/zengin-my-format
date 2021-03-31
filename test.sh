
# --------- module.exports= と、ケツカンマをなくす ------------
# reg2='s/^module.exports[^\S]?=[^\S]?//'
# reg3="^module.exports\s?=\s?|;\n?$"
a="module.exports = {};;"
echo $a
aiueo=`echo $a | sed 's/^module.exports[ \f\n\r\t]=[ \f\n\r\t]//' | sed 's/;$//'`
echo $aiueo

zengin=`curl https://raw.githubusercontent.com/zengin-code/zengin-js/master/lib/zengin-data.js`

reg1='s/^module.exports[ \f\n\r\t]=[ \f\n\r\t]//'
reg2='s/;$//'