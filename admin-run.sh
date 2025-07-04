#!/bin/bash

# Función que se ejecutará al presionar Ctrl+C
cleanup() {
    echo -e "\n\nCerrando puerto 8081..."
    sudo ufw deny 8081
    sudo ufw reload
    echo "Puerto 8081 cerrado."
    exit 0
}

# Configurar trap para capturar Ctrl+C (SIGINT)
trap cleanup INT

# Abrir puerto 8081
echo "Abriendo puerto 8081..."
sudo ufw allow 8081
sudo ufw --force enable

# Ejecutar air
echo "Iniciando air..."
air -c .air.toml  # Ajusta los parámetros según necesites

# Esta línea nunca se ejecutará si se usa Ctrl+C
cleanup