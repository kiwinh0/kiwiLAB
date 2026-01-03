#!/usr/bin/env python3
import json
import os

i18n_dir = '/Users/kiwinho/Proyectos/CodigoSH/web/static/i18n'

# Translations for login page
login_translations = {
    'en': {
        'login': {
            'tagline': 'Your digital command center',
            'username_label': 'Username',
            'username_placeholder': 'Enter your username',
            'password_label': 'Password',
            'password_placeholder': 'Enter your password',
            'remember_me': 'Keep me signed in',
            'remember_hint': 'Active session for 30 days',
            'submit_button': 'Sign In'
        }
    },
    'es': {
        'login': {
            'tagline': 'Tu centro de comando digital',
            'username_label': 'Usuario',
            'username_placeholder': 'Ingresa tu usuario',
            'password_label': 'Contraseña',
            'password_placeholder': 'Ingresa tu contraseña',
            'remember_me': 'Mantenerme conectado',
            'remember_hint': 'Sesión activa por 30 días',
            'submit_button': 'Iniciar Sesión'
        }
    },
    'fr': {
        'login': {
            'tagline': 'Votre centre de commandement numérique',
            'username_label': 'Nom d\'utilisateur',
            'username_placeholder': 'Entrez votre nom d\'utilisateur',
            'password_label': 'Mot de passe',
            'password_placeholder': 'Entrez votre mot de passe',
            'remember_me': 'Me garder connecté',
            'remember_hint': 'Session active pendant 30 jours',
            'submit_button': 'Se connecter'
        }
    },
    'de': {
        'login': {
            'tagline': 'Dein digitales Kommandozentrum',
            'username_label': 'Benutzername',
            'username_placeholder': 'Geben Sie Ihren Benutzernamen ein',
            'password_label': 'Passwort',
            'password_placeholder': 'Geben Sie Ihr Passwort ein',
            'remember_me': 'Anmeldedaten speichern',
            'remember_hint': 'Aktive Sitzung für 30 Tage',
            'submit_button': 'Anmelden'
        }
    },
    'it': {
        'login': {
            'tagline': 'Il tuo centro di comando digitale',
            'username_label': 'Nome utente',
            'username_placeholder': 'Inserisci il tuo nome utente',
            'password_label': 'Password',
            'password_placeholder': 'Inserisci la tua password',
            'remember_me': 'Mantienimi connesso',
            'remember_hint': 'Sessione attiva per 30 giorni',
            'submit_button': 'Accedi'
        }
    },
    'pt': {
        'login': {
            'tagline': 'Seu centro de comando digital',
            'username_label': 'Nome de usuário',
            'username_placeholder': 'Digite seu nome de usuário',
            'password_label': 'Senha',
            'password_placeholder': 'Digite sua senha',
            'remember_me': 'Mantenha-me conectado',
            'remember_hint': 'Sessão ativa por 30 dias',
            'submit_button': 'Entrar'
        }
    },
    'ru': {
        'login': {
            'tagline': 'Ваш центр цифрового управления',
            'username_label': 'Имя пользователя',
            'username_placeholder': 'Введите ваше имя пользователя',
            'password_label': 'Пароль',
            'password_placeholder': 'Введите ваш пароль',
            'remember_me': 'Держать меня в системе',
            'remember_hint': 'Активная сессия на 30 дней',
            'submit_button': 'Войти'
        }
    },
    'zh': {
        'login': {
            'tagline': '您的数字指挥中心',
            'username_label': '用户名',
            'username_placeholder': '输入您的用户名',
            'password_label': '密码',
            'password_placeholder': '输入您的密码',
            'remember_me': '保持登录',
            'remember_hint': '30天内保持活跃会话',
            'submit_button': '登录'
        }
    },
    'ja': {
        'login': {
            'tagline': 'あなたのデジタルコマンドセンター',
            'username_label': 'ユーザー名',
            'username_placeholder': 'ユーザー名を入力',
            'password_label': 'パスワード',
            'password_placeholder': 'パスワードを入力',
            'remember_me': 'ログインしたままにする',
            'remember_hint': '30日間有効なセッション',
            'submit_button': 'サインイン'
        }
    },
    'ko': {
        'login': {
            'tagline': '당신의 디지털 명령 센터',
            'username_label': '사용자명',
            'username_placeholder': '사용자명을 입력하세요',
            'password_label': '비밀번호',
            'password_placeholder': '비밀번호를 입력하세요',
            'remember_me': '로그인 유지',
            'remember_hint': '30일 동안 활성 세션',
            'submit_button': '로그인'
        }
    },
    'ar': {
        'login': {
            'tagline': 'مركز القيادة الرقمي الخاص بك',
            'username_label': 'اسم المستخدم',
            'username_placeholder': 'أدخل اسم المستخدم الخاص بك',
            'password_label': 'كلمة المرور',
            'password_placeholder': 'أدخل كلمة المرور الخاصة بك',
            'remember_me': 'ابقني مسجل الدخول',
            'remember_hint': 'جلسة نشطة لمدة 30 يومًا',
            'submit_button': 'دخول'
        }
    },
    'hi': {
        'login': {
            'tagline': 'आपका डिजिटल कमांड सेंटर',
            'username_label': 'उपयोगकर्ता नाम',
            'username_placeholder': 'अपना उपयोगकर्ता नाम दर्ज करें',
            'password_label': 'पासवर्ड',
            'password_placeholder': 'अपना पासवर्ड दर्ज करें',
            'remember_me': 'मुझे लॉगिन रखें',
            'remember_hint': '30 दिनों के लिए सक्रिय सेशन',
            'submit_button': 'साइन इन करें'
        }
    },
    'nl': {
        'login': {
            'tagline': 'Uw digitale commandocentrum',
            'username_label': 'Gebruikersnaam',
            'username_placeholder': 'Voer uw gebruikersnaam in',
            'password_label': 'Wachtwoord',
            'password_placeholder': 'Voer uw wachtwoord in',
            'remember_me': 'Houd me aangemeld',
            'remember_hint': 'Actieve sessie gedurende 30 dagen',
            'submit_button': 'Aanmelden'
        }
    },
    'sv': {
        'login': {
            'tagline': 'Ditt digitala kommandocenter',
            'username_label': 'Användarnamn',
            'username_placeholder': 'Ange ditt användarnamn',
            'password_label': 'Lösenord',
            'password_placeholder': 'Ange ditt lösenord',
            'remember_me': 'Håll mig inloggad',
            'remember_hint': 'Aktiv session i 30 dagar',
            'submit_button': 'Logga in'
        }
    },
    'pl': {
        'login': {
            'tagline': 'Twoje cyfrowe centrum dowodzenia',
            'username_label': 'Nazwa użytkownika',
            'username_placeholder': 'Wpisz swoją nazwę użytkownika',
            'password_label': 'Hasło',
            'password_placeholder': 'Wpisz swoje hasło',
            'remember_me': 'Zapamiętaj mnie',
            'remember_hint': 'Aktywna sesja przez 30 dni',
            'submit_button': 'Zaloguj się'
        }
    },
    'tr': {
        'login': {
            'tagline': 'Dijital komuta merkeziniz',
            'username_label': 'Kullanıcı Adı',
            'username_placeholder': 'Kullanıcı adınızı girin',
            'password_label': 'Şifre',
            'password_placeholder': 'Şifrenizi girin',
            'remember_me': 'Beni oturum açık tut',
            'remember_hint': '30 gün boyunca etkin oturum',
            'submit_button': 'Giriş Yap'
        }
    },
    'vi': {
        'login': {
            'tagline': 'Trung tâm chỉ huy kỹ thuật số của bạn',
            'username_label': 'Tên đăng nhập',
            'username_placeholder': 'Nhập tên đăng nhập của bạn',
            'password_label': 'Mật khẩu',
            'password_placeholder': 'Nhập mật khẩu của bạn',
            'remember_me': 'Giữ tôi đăng nhập',
            'remember_hint': 'Phiên hoạt động trong 30 ngày',
            'submit_button': 'Đăng nhập'
        }
    },
    'th': {
        'login': {
            'tagline': 'ศูนย์ควบคุมดิจิทัลของคุณ',
            'username_label': 'ชื่อผู้ใช้',
            'username_placeholder': 'ป้อนชื่อผู้ใช้ของคุณ',
            'password_label': 'รหัสผ่าน',
            'password_placeholder': 'ป้อนรหัสผ่านของคุณ',
            'remember_me': 'ให้ฉันอยู่ในสถานะเข้าสู่ระบบ',
            'remember_hint': 'เซสชั่นที่ใช้งานอยู่เป็นเวลา 30 วัน',
            'submit_button': 'ลงชื่อเข้า'
        }
    },
    'id': {
        'login': {
            'tagline': 'Pusat perintah digital Anda',
            'username_label': 'Nama pengguna',
            'username_placeholder': 'Masukkan nama pengguna Anda',
            'password_label': 'Kata sandi',
            'password_placeholder': 'Masukkan kata sandi Anda',
            'remember_me': 'Biarkan saya tetap masuk',
            'remember_hint': 'Sesi aktif selama 30 hari',
            'submit_button': 'Masuk'
        }
    },
    'el': {
        'login': {
            'tagline': 'Το κέντρο ψηφιακής διοίκησής σας',
            'username_label': 'Όνομα χρήστη',
            'username_placeholder': 'Εισάγετε το όνομα χρήστη σας',
            'password_label': 'Κωδικός πρόσβασης',
            'password_placeholder': 'Εισάγετε τον κωδικό πρόσβασής σας',
            'remember_me': 'Διατήρηση σύνδεσης',
            'remember_hint': 'Ενεργή σύνδεση για 30 ημέρες',
            'submit_button': 'Σύνδεση'
        }
    }
}

# Translations for about page
about_translations = {
    'en': {
        'about': {
            'version': 'Version 1.0.0',
            'description': 'A modern web application to manage your bookmarks and favorite links. Organize, customize and quickly access your most important resources.',
            'features': 'Main Features',
            'developed_with': 'Developed with ❤️ using Go and modern web technologies',
            'copyright': '© 2026 CodigoSH. All rights reserved.',
            'view_on_github': 'View on GitHub'
        }
    },
    'es': {
        'about': {
            'version': 'Versión 1.0.0',
            'description': 'Una aplicación web moderna para gestionar tus marcadores y enlaces favoritos. Organiza, personaliza y accede rápidamente a tus recursos más importantes.',
            'features': 'Características Principales',
            'developed_with': 'Desarrollado con ❤️ usando Go y modern web technologies',
            'copyright': '© 2026 CodigoSH. Todos los derechos reservados.',
            'view_on_github': 'Ver en GitHub'
        }
    },
    'fr': {
        'about': {
            'version': 'Version 1.0.0',
            'description': 'Une application web moderne pour gérer vos signets et liens préférés. Organisez, personnalisez et accédez rapidement à vos ressources les plus importantes.',
            'features': 'Fonctionnalités Principales',
            'developed_with': 'Développé avec ❤️ en utilisant Go et les technologies web modernes',
            'copyright': '© 2026 CodigoSH. Tous droits réservés.',
            'view_on_github': 'Voir sur GitHub'
        }
    },
    'de': {
        'about': {
            'version': 'Version 1.0.0',
            'description': 'Eine moderne Webanwendung zur Verwaltung Ihrer Lesezeichen und Lieblingslinks. Organisieren, personalisieren Sie und greifen Sie schnell auf Ihre wichtigsten Ressourcen zu.',
            'features': 'Hauptmerkmale',
            'developed_with': 'Entwickelt mit ❤️ unter Verwendung von Go und modernen Webtechnologien',
            'copyright': '© 2026 CodigoSH. Alle Rechte vorbehalten.',
            'view_on_github': 'Auf GitHub anzeigen'
        }
    },
    'it': {
        'about': {
            'version': 'Versione 1.0.0',
            'description': 'Un\'applicazione web moderna per gestire i tuoi segnalibri e link preferiti. Organizza, personalizza e accedi rapidamente alle tue risorse più importanti.',
            'features': 'Caratteristiche Principali',
            'developed_with': 'Sviluppato con ❤️ utilizzando Go e tecnologie web moderne',
            'copyright': '© 2026 CodigoSH. Tutti i diritti riservati.',
            'view_on_github': 'Visualizza su GitHub'
        }
    },
    'pt': {
        'about': {
            'version': 'Versão 1.0.0',
            'description': 'Uma aplicação web moderna para gerenciar seus marcadores e links favoritos. Organize, personalize e acesse rapidamente seus recursos mais importantes.',
            'features': 'Principais Características',
            'developed_with': 'Desenvolvido com ❤️ usando Go e tecnologias web modernas',
            'copyright': '© 2026 CodigoSH. Todos os direitos reservados.',
            'view_on_github': 'Ver no GitHub'
        }
    },
    'ru': {
        'about': {
            'version': 'Версия 1.0.0',
            'description': 'Современное веб-приложение для управления закладками и избранными ссылками. Организуйте, персонализируйте и быстро получайте доступ к своим наиболее важным ресурсам.',
            'features': 'Основные Возможности',
            'developed_with': 'Разработано с ❤️ с использованием Go и современных веб-технологий',
            'copyright': '© 2026 CodigoSH. Все права защищены.',
            'view_on_github': 'Просмотреть на GitHub'
        }
    },
    'zh': {
        'about': {
            'version': '版本 1.0.0',
            'description': '一个现代网络应用程序，用于管理您的书签和最喜欢的链接。组织、自定义并快速访问您最重要的资源。',
            'features': '主要功能',
            'developed_with': '用 ❤️ 使用 Go 和现代网络技术开发',
            'copyright': '© 2026 CodigoSH。版权所有。',
            'view_on_github': '在 GitHub 上查看'
        }
    },
    'ja': {
        'about': {
            'version': 'バージョン 1.0.0',
            'description': 'ブックマークとお気に入りのリンクを管理するための最新のウェブアプリケーション。整理し、カスタマイズし、最も重要なリソースにすばやくアクセスします。',
            'features': '主な機能',
            'developed_with': '❤️ を使用して Go と最新の Web テクノロジーで開発',
            'copyright': '© 2026 CodigoSH。著作権所有。',
            'view_on_github': 'GitHub で表示'
        }
    },
    'ko': {
        'about': {
            'version': '버전 1.0.0',
            'description': '북마크 및 선호하는 링크를 관리하기 위한 최신 웹 애플리케이션입니다. 정리하고, 사용자 정의하고, 가장 중요한 리소스에 빠르게 액세스하세요.',
            'features': '주요 기능',
            'developed_with': '❤️ 로 Go 및 최신 웹 기술을 사용하여 개발됨',
            'copyright': '© 2026 CodigoSH. 모든 권리 보유.',
            'view_on_github': 'GitHub에서 보기'
        }
    },
    'ar': {
        'about': {
            'version': 'الإصدار 1.0.0',
            'description': 'تطبيق ويب حديث لإدارة الإشارات المرجعية والروابط المفضلة لديك. تنظيم وتخصيص والوصول السريع إلى مواردك الأكثر أهمية.',
            'features': 'الميزات الرئيسية',
            'developed_with': 'تم تطويره بـ ❤️ باستخدام Go وتقنيات الويب الحديثة',
            'copyright': '© 2026 CodigoSH. جميع الحقوق محفوظة.',
            'view_on_github': 'عرض على GitHub'
        }
    },
    'hi': {
        'about': {
            'version': 'संस्करण 1.0.0',
            'description': 'अपने बुकमार्क और पसंदीदा लिंक का प्रबंधन करने के लिए एक आधुनिक वेब एप्लिकेशन। व्यवस्थित करें, कस्टमाइज़ करें और अपने सबसे महत्वपूर्ण संसाधनों तक तेज़ी से पहुंचें।',
            'features': 'मुख्य विशेषताएं',
            'developed_with': '❤️ के साथ Go और आधुनिक वेब प्रौद्योगिकी का उपयोग करके विकसित',
            'copyright': '© 2026 CodigoSH. सभी अधिकार सुरक्षित।',
            'view_on_github': 'GitHub पर देखें'
        }
    },
    'nl': {
        'about': {
            'version': 'Versie 1.0.0',
            'description': 'Een moderne webtoepassing voor het beheren van uw bladwijzers en favoriete koppelingen. Organiseer, personaliseer en krijg snel toegang tot uw belangrijkste bronnen.',
            'features': 'Hoofdfuncties',
            'developed_with': 'Ontwikkeld met ❤️ met Go en moderne webtechnologieën',
            'copyright': '© 2026 CodigoSH. Alle rechten voorbehouden.',
            'view_on_github': 'Weergeven op GitHub'
        }
    },
    'sv': {
        'about': {
            'version': 'Version 1.0.0',
            'description': 'En modern webbapplikation för att hantera dina bokmärken och favoritlänkar. Organisera, anpassa och få snabb åtkomst till dina viktigaste resurser.',
            'features': 'Huvudfunktioner',
            'developed_with': 'Utvecklad med ❤️ med Go och moderna webteknologier',
            'copyright': '© 2026 CodigoSH. Alla rättigheter förbehålles.',
            'view_on_github': 'Visa på GitHub'
        }
    },
    'pl': {
        'about': {
            'version': 'Wersja 1.0.0',
            'description': 'Nowoczesna aplikacja internetowa do zarządzania zakładkami i ulubionymi linkami. Organizuj, dostosowuj i szybko uzyskiwaj dostęp do swoich najważniejszych zasobów.',
            'features': 'Główne Funkcje',
            'developed_with': 'Opracowano z ❤️ przy użyciu Go i nowoczesnych technologii internetowych',
            'copyright': '© 2026 CodigoSH. Wszystkie prawa zastrzeżone.',
            'view_on_github': 'Wyświetl na GitHub'
        }
    },
    'tr': {
        'about': {
            'version': 'Sürüm 1.0.0',
            'description': 'Yer İşaretlerinizi ve favori bağlantılarınızı yönetmek için modern bir web uygulaması. En önemli kaynaklarınızı düzenleyin, özelleştirin ve hızlı bir şekilde erişin.',
            'features': 'Ana Özellikler',
            'developed_with': '❤️ ile Go ve modern web teknolojileri kullanılarak geliştirilmiştir',
            'copyright': '© 2026 CodigoSH. Tüm hakları saklıdır.',
            'view_on_github': 'GitHub\'da Görüntüle'
        }
    },
    'vi': {
        'about': {
            'version': 'Phiên bản 1.0.0',
            'description': 'Một ứng dụng web hiện đại để quản lý các dấu trang và liên kết yêu thích của bạn. Tổ chức, tùy chỉnh và truy cập nhanh các tài nguyên quan trọng nhất của bạn.',
            'features': 'Các Tính Năng Chính',
            'developed_with': 'Được phát triển với ❤️ bằng Go và công nghệ web hiện đại',
            'copyright': '© 2026 CodigoSH. Bảo lưu mọi quyền.',
            'view_on_github': 'Xem trên GitHub'
        }
    },
    'th': {
        'about': {
            'version': 'เวอร์ชั่น 1.0.0',
            'description': 'แอปพลิเคชันเว็บสมัยใหม่เพื่อจัดการที่คั่นหน้าและลิงก์ที่ชื่นชอบของคุณ จัดระเบียบ ปรับแต่ง และเข้าถึงทรัพยากรที่สำคัญที่สุดของคุณอย่างรวดเร็ว',
            'features': 'คุณสมบัติหลัก',
            'developed_with': 'พัฒนาด้วย ❤️ โดยใช้ Go และเทคโนโลยีเว็บสมัยใหม่',
            'copyright': '© 2026 CodigoSH. สงวนสิทธิ์ทั้งหมด',
            'view_on_github': 'ดูบน GitHub'
        }
    },
    'id': {
        'about': {
            'version': 'Versi 1.0.0',
            'description': 'Aplikasi web modern untuk mengelola penanda halaman dan tautan favorit Anda. Atur, sesuaikan, dan akses dengan cepat sumber daya paling penting Anda.',
            'features': 'Fitur Utama',
            'developed_with': 'Dikembangkan dengan ❤️ menggunakan Go dan teknologi web modern',
            'copyright': '© 2026 CodigoSH. Semua hak dilindungi.',
            'view_on_github': 'Lihat di GitHub'
        }
    },
    'el': {
        'about': {
            'version': 'Έκδοση 1.0.0',
            'description': 'Μια σύγχρονη εφαρμογή ιστού για τη διαχείριση των σελιδοδεικτών και των αγαπημένων συνδέσμων σας. Οργανώστε, προσαρμόστε και αποκτήστε γρήγορη πρόσβαση στους σημαντικότερους πόρους σας.',
            'features': 'Κύριες Λειτουργίες',
            'developed_with': 'Αναπτύχθηκε με ❤️ χρησιμοποιώντας Go και σύγχρονες τεχνολογίες ιστού',
            'copyright': '© 2026 CodigoSH. Όλα τα δικαιώματα διατηρούνται.',
            'view_on_github': 'Προβολή στο GitHub'
        }
    }
}

# Update JSON files
languages = ['en', 'es', 'fr', 'de', 'it', 'pt', 'ru', 'zh', 'ja', 'ko', 'ar', 'hi', 'nl', 'sv', 'pl', 'tr', 'vi', 'th', 'id', 'el']

for lang in languages:
    file_path = os.path.join(i18n_dir, f'{lang}.json')
    
    with open(file_path, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    # Update login translations
    if lang in login_translations:
        if 'login' not in data:
            data['login'] = {}
        data['login'].update(login_translations[lang]['login'])
    
    # Update about translations
    if lang in about_translations:
        if 'about' not in data:
            data['about'] = {}
        data['about'].update(about_translations[lang]['about'])
    
    with open(file_path, 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=4)
    
    print(f'✓ Updated {lang}.json with login and about translations')

print('✅ All translations added successfully!')
