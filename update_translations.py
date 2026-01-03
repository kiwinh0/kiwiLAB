#!/usr/bin/env python3
import json
import os
from pathlib import Path

i18n_dir = Path("/Users/kiwinho/Proyectos/CodigoSH/web/static/i18n")

# Translations for 'common' section in different languages
common_translations = {
    'en': {
        'save': 'Save Preferences',
        'cancel': 'Cancel',
        'enabled': 'Enabled',
        'disabled': 'Disabled'
    },
    'es': {
        'save': 'Guardar preferencias',
        'cancel': 'Cancelar',
        'enabled': 'Activado',
        'disabled': 'Desactivado'
    },
    'fr': {
        'save': 'Enregistrer les préférences',
        'cancel': 'Annuler',
        'enabled': 'Activé',
        'disabled': 'Désactivé'
    },
    'de': {
        'save': 'Einstellungen speichern',
        'cancel': 'Abbrechen',
        'enabled': 'Aktiviert',
        'disabled': 'Deaktiviert'
    },
    'it': {
        'save': 'Salva preferenze',
        'cancel': 'Annulla',
        'enabled': 'Abilitato',
        'disabled': 'Disabilitato'
    },
    'pt': {
        'save': 'Salvar preferências',
        'cancel': 'Cancelar',
        'enabled': 'Ativado',
        'disabled': 'Desativado'
    },
    'ru': {
        'save': 'Сохранить предпочтения',
        'cancel': 'Отмена',
        'enabled': 'Включено',
        'disabled': 'Отключено'
    },
    'zh': {
        'save': '保存首选项',
        'cancel': '取消',
        'enabled': '已启用',
        'disabled': '已禁用'
    },
    'ja': {
        'save': '設定を保存',
        'cancel': 'キャンセル',
        'enabled': '有効',
        'disabled': '無効'
    },
    'ko': {
        'save': '기본 설정 저장',
        'cancel': '취소',
        'enabled': '활성화됨',
        'disabled': '비활성화됨'
    },
    'ar': {
        'save': 'حفظ التفضيلات',
        'cancel': 'إلغاء',
        'enabled': 'مفعل',
        'disabled': 'معطل'
    },
    'hi': {
        'save': 'वरीयताएं सहेजें',
        'cancel': 'रद्द करें',
        'enabled': 'सक्षम',
        'disabled': 'अक्षम'
    },
    'nl': {
        'save': 'Voorkeuren opslaan',
        'cancel': 'Annuleren',
        'enabled': 'Ingeschakeld',
        'disabled': 'Uitgeschakeld'
    },
    'sv': {
        'save': 'Spara inställningar',
        'cancel': 'Avbryt',
        'enabled': 'Aktiverad',
        'disabled': 'Inaktiverad'
    },
    'pl': {
        'save': 'Zapisz preferencje',
        'cancel': 'Anuluj',
        'enabled': 'Włączone',
        'disabled': 'Wyłączone'
    },
    'tr': {
        'save': 'Tercihleri Kaydet',
        'cancel': 'İptal',
        'enabled': 'Etkin',
        'disabled': 'Devre Dışı'
    },
    'vi': {
        'save': 'Lưu tùy chọn',
        'cancel': 'Hủy',
        'enabled': 'Bật',
        'disabled': 'Tắt'
    },
    'th': {
        'save': 'บันทึกการตั้งค่า',
        'cancel': 'ยกเลิก',
        'enabled': 'เปิดใช้งาน',
        'disabled': 'ปิดใช้งาน'
    },
    'id': {
        'save': 'Simpan Preferensi',
        'cancel': 'Batal',
        'enabled': 'Diaktifkan',
        'disabled': 'Dinonaktifkan'
    },
    'el': {
        'save': 'Αποθήκευση προτιμήσεων',
        'cancel': 'Ακύρωση',
        'enabled': 'Ενεργοποιημένο',
        'disabled': 'Απενεργοποιημένο'
    }
}

# Process each JSON file
for json_file in i18n_dir.glob('*.json'):
    lang_code = json_file.stem
    
    with open(json_file, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    # Add common section if not exists or update it
    if lang_code in common_translations:
        data['common'] = common_translations[lang_code]
        print(f"✓ Updated {json_file.name}")
    
    with open(json_file, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=4)

print("\n✅ All translation files updated!")
