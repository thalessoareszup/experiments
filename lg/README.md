# lg

Lê JSON do stdin e renderiza numa TUI, util para explorar o output de um processo que emite JSON.

Você pode instalar com

```
go build -o lg main.go && mv lg ~/.local/bin
```

Assumindo que `~/.local/bin` esteja no seu `PATH`.

Depois, use com

```
kubectl logs my-pod-abc234 -f | lg
```
