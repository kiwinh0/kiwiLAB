import json
import glob

# Traducciones para "Añadir sección"
translations = {
    "es": "Añadir sección",
    "en": "Add section",
    "fr": "Ajouter une section",
    "de": "Abschnitt hinzufügen",
    "it": "Aggiungi sezione",
    "pt": "Adicionar seção",
    "ru": "Добавить раздел",
    "zh": "添加部分",
    "ja": "セクションを追加",
    "ko": "섹션 추가",
    "ar": "إضافة قسم",
    "hi": "अनुभाग जोड़ें",
    "nl": "Sectie toevoegen",
    "sv": "Lägg till sektion",
    "pl": "Dodaj sekcję",
    "tr": "Bölüm ekle",
    "vi": "Thêm phần",
    "th": "เพิ่มส่วน",
    "id": "Tambah bagian",
    "el": "Προσθήκη ενότητας"
}

for json_file in glob.glob("*.json"):
    lang = json_file.replace(".json", "")
    
    with open(json_file, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    if "dashboard" in data and "addSection" not in data["dashboard"]:
        data["dashboard"]["addSection"] = translations.get(lang, "Add section")
        
        with open(json_file, 'w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=4)
        
        print(f"✅ {json_file}: Añadido 'addSection'")
    else:
        print(f"⏭️  {json_file}: Ya existe 'addSection' o no tiene dashboard")

print("\n✅ Proceso completado")
