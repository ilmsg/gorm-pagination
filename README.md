# gorm-pagination
ตัวอย่างการทำ pagination กับ gorm

ชอบ func ตัวนี้

    func (*gorm.DB).Scopes(funcs ...func(*gorm.DB) *gorm.DB) (tx *gorm.DB)

เราส่งผ่าน `func(DB) DB` ที่เชื่อมต่อไว้แล้ว เข้าไปเพื่อใช้ กับ เงื่อนไข ที่เราต้องการได้

    func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
        return db.Where("amount > ?", 1000)
    }

    func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
        return func (db *gorm.DB) *gorm.DB {
            return db.Scopes(AmountGreaterThan1000).Where("status in (?)", status)
        }
    }

    db.Scopes(
        AmountGreaterThan1000,
        OrderStatus([]string{"paid", "shipped"}),
    ).Find(&orders)

ตัวอย่างตอน

    $ go run . <limit> <page>


อย่างเช่น ต้องการหน้าล่ะ 5 รายการ, เปิดดูหน้า 2

    $ go run 5 2


ไปเจอตัวอย่างอันนี้มา, ในคอมเม้น มีคนบอกให้ระวังเรื่อง sql injection ตรง sort ด้วยนะ, ไม่ควรรับอะไรมาแล้วเอามาใช้ตรงๆแบบนั้น

refer: 

gorm pagination
https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

gorm seeding
https://ieftimov.com/posts/simple-golang-database-seeding-abstraction-gorm/