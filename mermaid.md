```mermaid 
sequenceDiagram
    participant Product
    participant Category

    
```

```mermaid
classDiagram
    class Product {
        +ID
        +Code
        +Name
        +Price
        +CategoryID
    }
    class Category {
        +ID
        +Code
        +Name
    }
    Product --> Category : belongs to
```