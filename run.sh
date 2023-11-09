#!bin/bash

echo -n "Enter the name of an animal: "
read ANIMAL
echo -n "The $ANIMAL has "
case $ANIMAL in
  horse | dog | cat) echo -n "four";;
  man | kangaroo ) echo -n "two";;
  *) echo -n "an unknown number of";;
esac
echo " legs."

# ORDER=$1
# echo "order=$ORDER"
# case "$ORDER" in
#     gen )
#         echo "order: $ORDER"
#     ;;
#
#     protoc )
#         echo "order: $ORDER"
#     ;;
#
#     *)
#         echo "invalid order: order=$ORDER"
#     ;;
# esac