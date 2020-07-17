#!/usr/bin/env bash

echo "***** deploy.sh -- Inicio ejecución -- ******"
echo ""
rm -rf C:/xampp/htdocs/nature/api
mkdir C:/xampp/htdocs/nature/api

echo " **** Copiando contenido de dist/ a htdocs/nature/api ****"
echo ""
cp -r ./dist/* C:/xampp/htdocs/nature/api
# cp -r ./dist/.* C:/xampp/htdocs/nature/api

echo ""
echo " ***** deploy.sh -- Fin ejecución -- ******"