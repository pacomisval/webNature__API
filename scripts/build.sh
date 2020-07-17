#!/usr/bin/env bash

echo " ****** build.sh -- Inicio ejecución -- *******"
echo ""
echo " ****** Borrando y Creando dist/ ******"
rm -rf ./dist/
mkdir ./dist/

echo " ***** Copiando contenido de src/ a dist/ ****"
echo ""
cp -r ./src/* ./dist/
echo ""
echo " ***** build.sh -- Fin ejecución -- ******"