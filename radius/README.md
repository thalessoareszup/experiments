# radius

Experimento local do https://docs.radapp.io/

## Instalação

A instalação padrão do radapp exige sudo. Para evitar isso, precisamos definir a variável RADIUS_INSTALL_DIR para uma pasta que temos acesso. O script irá baixar
e instalar um programa chamado `rad`, então confirme que RADIUS_INSTALL_DIR aponte para seu PATH.

```bash
wget -q "https://raw.githubusercontent.com/radius-project/radius/a45fc2f14ad7c3e2956e756fec212cf0d1805ba4/deploy/install.sh" -O - | RADIUS_INSTALL_DIR=$HOME/.local/bin /bin/bash
```

### Preparando um cluster de k8s

Rode o script `setup_k8s.sh`, que utilizará o `kind` para subir um cluster local chamado `radius-sandbox`, usando a configuração em `kind-config.yaml`, que injeta o certificado do Zscaler no cluster para que consiga fazer *pull* das imagens. Confirme se o caminho do certificado do Zscaler na sua máquina é o correto. O script também atualizará o contexto do `kubectl` para apontar para esse cluster.

### Inicializando o radius

Rode `rad initialize` na pasta `todoapp`. Isso instalará o radius no cluster e configurará uma aplicação padrão. Acompanhe o setup com `kubectl get pods -n radius-system`.

### Rodando uma aplicação

Rode `rad run app.bicep` na pasta `todoapp` para realizar um deploy de uma aplicação de exemplo.
