# 📊 Sistema de Movimientos de Clientes - Simplificado

## 🎯 Concepto

Sistema manual de contabilidad para clientes donde:
- Los movimientos NO están relacionados con productos/ventas
- Se registran manualmente con los datos que tú proporciones
- Hay 2 tipos de movimientos: **DEUDA** y **ABONO**
- Balance = Σ(DEUDA) - Σ(ABONO)
- Puede quedar saldo a favor del cliente (balance negativo)

---

## 🔧 Endpoint Principal

### **POST /api/v1/customers/transactions**

Agregar uno o múltiples movimientos a un cliente.

**Request Body:**
```json
{
  "customer_id": 1,
  "transactions": [
    {
      "type": "DEUDA",
      "amount": 150000,
      "description": "Compra de chaqueta de cuero y pantalón",
      "date": "2024-11-20T00:00:00Z"  // Opcional
    },
    {
      "type": "ABONO",
      "amount": 50000,
      "description": "Abono inicial en efectivo",
      "payment_method_id": 4,  // Requerido para ABONO
      "date": "2024-11-20T00:00:00Z"  // Opcional
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Transactions added successfully",
  "data": [
    {
      "id": 1,
      "customer_id": 1,
      "type": "DEUDA",
      "amount": 150000,
      "description": "Compra de chaqueta de cuero y pantalón",
      "date": "2024-11-20T00:00:00Z"
    },
    {
      "id": 2,
      "customer_id": 1,
      "type": "ABONO",
      "amount": 50000,
      "description": "Abono inicial en efectivo",
      "payment_method_id": 4,
      "payment_method": {
        "id": 4,
        "name": "Efectivo"
      },
      "date": "2024-11-20T00:00:00Z"
    }
  ]
}
```

---

## 📋 Casos de Uso

### **Caso 1: Compra sin abono (a crédito)**

Cliente compra varios items pero no paga nada.

```bash
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "transactions": [
      {
        "type": "DEUDA",
        "amount": 100000,
        "description": "Chaqueta de cuero negra talla M"
      },
      {
        "type": "DEUDA",
        "amount": 80000,
        "description": "Pantalón jean azul talla 32"
      },
      {
        "type": "DEUDA",
        "amount": 90000,
        "description": "Zapatos deportivos talla 9.5"
      }
    ]
  }'
```

**Resultado:**
- Total deuda: 270,000
- Balance: +270,000

---

### **Caso 2: Compra con abono parcial**

Cliente compra y paga una parte.

```bash
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "transactions": [
      {
        "type": "DEUDA",
        "amount": 50000,
        "description": "Camiseta polo blanca"
      },
      {
        "type": "DEUDA",
        "amount": 80000,
        "description": "Pantalón formal negro"
      },
      {
        "type": "ABONO",
        "amount": 70000,
        "description": "Abono parcial - NEQUI",
        "payment_method_id": 1
      }
    ]
  }'
```

**Resultado:**
- Total deuda: 130,000
- Total abonado: 70,000
- Balance: +60,000

---

### **Caso 3: Compra con abono total (de esa compra)**

Cliente compra y paga exactamente el total de esa compra.

```bash
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "transactions": [
      {
        "type": "DEUDA",
        "amount": 50000,
        "description": "Gorra deportiva"
      },
      {
        "type": "DEUDA",
        "amount": 30000,
        "description": "Medias x3 pares"
      },
      {
        "type": "ABONO",
        "amount": 80000,
        "description": "Pago completo en efectivo",
        "payment_method_id": 4
      }
    ]
  }'
```

**Resultado:**
- Total deuda: 80,000
- Total abonado: 80,000
- Balance de esta compra: 0
- ⚠️ Si tenía deuda anterior, el balance total NO será 0

---

### **Caso 4: Compra con abono que cubre TODO el balance**

Cliente compra y paga TODO lo que debe (incluyendo deudas anteriores).

```bash
# Primero ver el balance actual
curl -X GET http://localhost:8080/api/v1/customers/1/balance \
  -H "Authorization: Bearer $TOKEN"
# Respuesta: { "balance": 150000 }

# Agregar nueva compra y pagar TODO
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "transactions": [
      {
        "type": "DEUDA",
        "amount": 50000,
        "description": "Chaqueta nueva"
      },
      {
        "type": "ABONO",
        "amount": 200000,
        "description": "Pago total del balance - Daviplata",
        "payment_method_id": 3
      }
    ]
  }'
```

**Resultado:**
- Balance anterior: 150,000
- Nueva deuda: +50,000
- Total a pagar: 200,000
- Abono: -200,000
- **Balance final: 0**

---

### **Caso 5: Solo abono (sin compra)**

Cliente viene solo a pagar, sin comprar nada.

```bash
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "transactions": [
      {
        "type": "ABONO",
        "amount": 50000,
        "description": "Abono quincenal - día 2",
        "payment_method_id": 2,
        "date": "2024-11-02T00:00:00Z"
      }
    ]
  }'
```

**Resultado:**
- No se agrega deuda
- Solo se registra el abono
- Balance disminuye en 50,000

---

### **Caso 6: Abono mayor que la deuda (saldo a favor)**

Cliente paga más de lo que debe, queda con saldo a favor.

```bash
# Balance actual: 80,000
curl -X POST http://localhost:8080/api/v1/customers/transactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "transactions": [
      {
        "type": "ABONO",
        "amount": 100000,
        "description": "Pago anticipado para próximas compras",
        "payment_method_id": 1
      }
    ]
  }'
```

**Resultado:**
- Balance anterior: +80,000
- Abono: -100,000
- **Balance final: -20,000 (saldo a favor)**

---

## 🔍 Consultar Balance

```bash
curl -X GET http://localhost:8080/api/v1/customers/1/balance \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta:**
```json
{
  "success": true,
  "message": "Balance retrieved successfully",
  "data": {
    "customer_id": 1,
    "balance": 100000
  }
}
```

**Interpretación del balance:**
- **Positivo (+100,000)**: Cliente debe 100,000
- **Cero (0)**: Cliente no debe nada
- **Negativo (-20,000)**: Cliente tiene saldo a favor de 20,000

---

## 📜 Ver Historial

```bash
curl -X GET http://localhost:8080/api/v1/customers/1/history \
  -H "Authorization: Bearer $TOKEN"
```

**Respuesta:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "customer_id": 1,
      "type": "DEUDA",
      "amount": 150000,
      "description": "Chaqueta de cuero y pantalón",
      "date": "2024-11-15T10:00:00Z"
    },
    {
      "id": 2,
      "customer_id": 1,
      "type": "ABONO",
      "amount": 50000,
      "description": "Abono quincenal",
      "payment_method_id": 1,
      "payment_method": {
        "id": 1,
        "name": "NEQUI Sonia"
      },
      "date": "2024-11-17T14:30:00Z"
    }
  ]
}
```

---

## 💡 Ventajas de este Sistema

1. ✅ **Flexibilidad Total**: Registras exactamente lo que vendiste
2. ✅ **Sin Dependencias**: No necesitas productos en el sistema
3. ✅ **Múltiples Movimientos**: Puedes agregar varios items y abonos en una sola llamada
4. ✅ **Saldo a Favor**: Soporta que el cliente pague de más
5. ✅ **Historial Completo**: Ves todos los movimientos con fechas y métodos de pago
6. ✅ **Balance Automático**: Se calcula dinámicamente

---

## 📊 Fórmula del Balance

```
Balance = Σ(DEUDA) - Σ(ABONO)
```

**Ejemplo:**
```
DEUDA:  +150,000 (chaqueta)
DEUDA:  +80,000 (pantalón)
ABONO:  -100,000 (pago)
DEUDA:  +50,000 (zapatos)
ABONO:  -50,000 (pago)
------------------------
Balance: +130,000 (debe)
```

---

## 🎨 Métodos de Pago Disponibles

| ID | Nombre |
|----|--------|
| 1  | NEQUI Sonia |
| 2  | NEQUI Jhon |
| 3  | Daviplata |
| 4  | Efectivo |

---

## ⚠️ Validaciones

1. **ABONO requiere payment_method_id**: No puedes registrar un abono sin especificar cómo pagó
2. **Amount siempre positivo**: El tipo (DEUDA/ABONO) define si suma o resta
3. **Customer debe existir**: El customer_id debe ser válido
4. **Al menos 1 transacción**: Debes enviar mínimo un movimiento
5. **Fecha opcional**: Si no envías fecha, se usa la actual

---

## 🚀 Próximos Pasos

1. Implementar en `main.go` y `routes.go`
2. Ejecutar migraciones para actualizar la tabla
3. Probar los casos de uso
4. Actualizar documentación API

¿Listo para continuar con la implementación?
