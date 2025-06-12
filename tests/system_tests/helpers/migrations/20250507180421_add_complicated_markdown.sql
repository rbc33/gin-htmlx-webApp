-- +goose Up
-- +goose StatementBegin
INSERT INTO posts(title, content, excerpt) 
VALUES(
    'My Very First Post',
    "# Permutación en Cadenas: Problema y Soluciones 🔄

## 📝 Descripción del Problema

Dados dos strings `s1` y `s2`, retornar `true` si `s2` contiene una permutación de `s1`, o `false` en caso contrario.

### Ejemplo
```
Input: s1 = 'ab', s2 = 'eidbaooo'
Output: true
Explicación: s2 contiene una permutación de s1 ('ba')
```

## 🧮 Fórmula Matemática

La probabilidad de encontrar una permutación en una cadena de longitud n:

$$ P(permutación) = \frac{n!}{(n-k)!} $$

Donde:
- n = longitud de s2
- k = longitud de s1

## 💡 Soluciones

### 1. Solución con Backtracking

```python
class Solution:
    def checkInclusion(self, s1: str, s2: str) -> bool:
        def backtrack(start, comb, visited):
            if len(comb) == len(s1):
                return comb in s2
            
            for i in range(len(s1)):
                if i in visited:
                    continue
                visited.append(i)
                new_comb = comb + s1[i]
                if backtrack(i, new_comb, visited):
                    return True
                visited.pop()
            return False
            
        return backtrack(0, '', [])
```

### 2. Solución con Ventana Deslizante

```python
class Solution:
    def checkInclusion(self, s1: str, s2: str) -> bool:
        if len(s1) > len(s2):
            return False
        
        s1_count = [0] * 26
        window_count = [0] * 26
        
        for char in s1:
            s1_count[ord(char) - ord('a')] += 1
        
        for i in range(len(s1)):
            window_count[ord(s2[i]) - ord('a')] += 1
        
        if s1_count == window_count:
            return True
        
        for i in range(len(s1), len(s2)):
            window_count[ord(s2[i]) - ord('a')] += 1
            window_count[ord(s2[i - len(s1)]) - ord('a')] -= 1
            
            if s1_count == window_count:
                return True
        
        return False
```

## 📊 Comparación de Complejidades

| Solución | Tiempo | Espacio | 
|----------|---------|---------|
| Backtracking | O(n!) | O(n) |
| Ventana Deslizante | O(n) | O(1) |

## 🔍 Análisis Detallado

### Ventana Deslizante
1. Inicialización
   - Crear arrays de conteo
   - Procesar primera ventana

2. Procesamiento
   ```mermaid
   graph LR
   A[Inicio Ventana] --> B[Contar Caracteres]
   B --> C[Deslizar Ventana]
   C --> D[Actualizar Conteos]
   D --> E[Comparar]
   ```

### Backtracking
1. Generación de permutaciones
   - Recursión
   - Control de visitados

## 📝 Casos de Prueba

| Input s1 | Input s2 | Output | Explicación |
|----------|----------|--------|-------------|
| 'ab' | 'eidbaooo' | true | Contiene 'ba' |
| 'ab' | 'eidboaoo' | false | No hay permutación |
| 'abc' | 'bbbca' | true | Contiene 'bca' |

## ⚠️ Consideraciones

1. Manejo de casos especiales:
   - Longitudes diferentes
   - Strings vacíos
   - Caracteres repetidos

2. Optimizaciones:
   ```python
   # Verificación rápida inicial
   if len(s1) > len(s2):
       return False
   ```

## 🔗 Referencias

- [Algoritmos de Backtracking](https://example.com)
- [Técnica de Ventana Deslizante](https://example.com)
- [Complejidad Algorítmica](https://example.com)

## 📊 Gráfico de Rendimiento

```ascii
Rendimiento vs Tamaño de Input
│
│    Ventana Deslizante
│    ┌─────────────
│    │
│    │      Backtracking
│    │      ╱
│    │     ╱
│    │    ╱
│    │   ╱
└────┴──────────────
```

## 🏁 Conclusión

La solución de ventana deslizante es superior en términos de:
- Eficiencia temporal
- Uso de memoria
- Claridad de código
- Mantenibilidad

---
*Nota: Este documento es parte de una serie de soluciones algorítmicas.*", "Una explicación detallada sobre el algoritmo de permutación en cadenas utilizando el método de ventana deslizante");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts ORDER BY id DESC LIMIT 1;

-- +goose StatementEnd
