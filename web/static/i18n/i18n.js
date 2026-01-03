// Sistema de internacionalización para CodigoSH
class I18n {
    constructor() {
        this.translations = {};
        this.supportedLanguages = [
            'en', 'es', 'fr', 'de', 'it', 'pt', 'ru', 'zh', 'ja', 'ko',
            'ar', 'hi', 'nl', 'sv', 'pl', 'tr', 'vi', 'th', 'id', 'el'
        ];
        this.currentLanguage = this.getStoredLanguage() || 'en';
    }

    getStoredLanguage() {
        // PRIORIDAD 1: HTML lang attribute (viene de la base de datos)
        const htmlLang = document.documentElement.getAttribute('lang');
        if (htmlLang && this.supportedLanguages.includes(htmlLang)) {
            // Sincronizar localStorage con la fuente de verdad (BD)
            const storedLang = localStorage.getItem('language');
            if (storedLang !== htmlLang) {
                console.log(`[i18n] Sincronizando localStorage: ${storedLang} → ${htmlLang}`);
                this.persistLanguage(htmlLang);
            }
            return htmlLang;
        }
        
        // PRIORIDAD 2: localStorage (solo si no hay HTML lang)
        const hasUserSelection = localStorage.getItem('languageSelected') === 'true';
        if (hasUserSelection) {
            return localStorage.getItem('language') || localStorage.getItem('selectedLanguage') || 'en';
        }
        
        // PRIORIDAD 3: Fallback a inglés
        return 'en';
    }

    async loadLanguage(lang) {
        if (!this.supportedLanguages.includes(lang)) {
            lang = 'en';
        }

        try {
            const response = await fetch(`/static/i18n/${lang}.json`);
            if (!response.ok) {
                throw new Error(`Failed to load language: ${lang}`);
            }
            this.translations = await response.json();
            this.currentLanguage = lang;
            // Do NOT save to localStorage automatically - only save when explicitly persisting
            this.updatePageLanguage();
            return true;
        } catch (error) {
            console.error('Error loading language:', error);
            // Fallback to English
            if (lang !== 'en') {
                return this.loadLanguage('en');
            }
            return false;
        }
    }

    persistLanguage(lang) {
        // Explicitly save language to localStorage
        localStorage.setItem('language', lang);
        localStorage.setItem('selectedLanguage', lang);
        localStorage.setItem('languageSelected', 'true');
        localStorage.setItem('currentLanguage', lang);
        
        // Also save to cookie for server-side access
        // Set cookie with 1 year expiration
        const expirationDate = new Date();
        expirationDate.setFullYear(expirationDate.getFullYear() + 1);
        document.cookie = `currentLanguage=${lang}; expires=${expirationDate.toUTCString()}; path=/`;
    }

    t(key) {
        const keys = key.split('.');
        let value = this.translations;
        
        for (const k of keys) {
            if (value && typeof value === 'object' && k in value) {
                value = value[k];
            } else {
                console.warn(`Translation key not found: ${key}`);
                return key;
            }
        }
        
        return value;
    }

    // Alias for t() for compatibility
    translate(key) {
        return this.t(key);
    }

    updatePageLanguage() {
        // Actualizar elementos con data-i18n
        document.querySelectorAll('[data-i18n]').forEach(element => {
            const key = element.getAttribute('data-i18n');
            const translation = this.t(key);
            
            if (element.tagName === 'INPUT' || element.tagName === 'TEXTAREA') {
                if (element.hasAttribute('placeholder')) {
                    element.placeholder = translation;
                } else {
                    element.value = translation;
                }
            } else {
                element.textContent = translation;
            }
        });

        // Actualizar placeholders con data-i18n-placeholder
        document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
            const key = element.getAttribute('data-i18n-placeholder');
            element.placeholder = this.t(key);
        });

        // Actualizar títulos con data-i18n-title
        document.querySelectorAll('[data-i18n-title]').forEach(element => {
            const key = element.getAttribute('data-i18n-title');
            element.title = this.t(key);
        });

        // Disparar evento personalizado
        document.dispatchEvent(new CustomEvent('languageChanged', {
            detail: { language: this.currentLanguage }
        }));
    }

    async setLanguage(lang) {
        await this.loadLanguage(lang);
    }

    getCurrentLanguage() {
        return this.currentLanguage;
    }

    getSupportedLanguages() {
        return this.supportedLanguages;
    }

    getLanguageName(code) {
        return this.t(`languages.${code}`);
    }
}

// Instancia global
const i18n = new I18n();
// Exponer en window para que otras partes (setup) puedan usarla
if (typeof window !== 'undefined') {
    window.i18n = i18n;
}

// Inicializar automáticamente cuando el DOM esté listo
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
        i18n.loadLanguage(i18n.getCurrentLanguage());
    });
} else {
    i18n.loadLanguage(i18n.getCurrentLanguage());
}
