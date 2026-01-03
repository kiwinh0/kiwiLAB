import json
import os

# Traducciones para los mensajes de alerta
alert_translations = {
    'es': {
        'profileUpdated': 'Perfil actualizado correctamente',
        'invalidPassword': 'Contraseña actual incorrecta',
        'passwordRequired': 'Debes ingresar tu contraseña actual para cambiar la contraseña',
        'passwordMismatch': 'Las contraseñas nuevas no coinciden',
        'usernameExists': 'El nombre de usuario ya existe',
        'invalidFileType': 'Tipo de archivo no válido. Solo se permiten imágenes.'
    },
    'en': {
        'profileUpdated': 'Profile updated successfully',
        'invalidPassword': 'Invalid current password',
        'passwordRequired': 'You must enter your current password to change it',
        'passwordMismatch': 'New passwords do not match',
        'usernameExists': 'Username already exists',
        'invalidFileType': 'Invalid file type. Only images are allowed.'
    },
    'fr': {
        'profileUpdated': 'Profil mis à jour avec succès',
        'invalidPassword': 'Mot de passe actuel incorrect',
        'passwordRequired': 'Vous devez saisir votre mot de passe actuel pour le modifier',
        'passwordMismatch': 'Les nouveaux mots de passe ne correspondent pas',
        'usernameExists': "Le nom d'utilisateur existe déjà",
        'invalidFileType': "Type de fichier invalide. Seules les images sont autorisées."
    },
    'de': {
        'profileUpdated': 'Profil erfolgreich aktualisiert',
        'invalidPassword': 'Ungültiges aktuelles Passwort',
        'passwordRequired': 'Sie müssen Ihr aktuelles Passwort eingeben, um es zu ändern',
        'passwordMismatch': 'Neue Passwörter stimmen nicht überein',
        'usernameExists': 'Benutzername existiert bereits',
        'invalidFileType': 'Ungültiger Dateityp. Nur Bilder sind erlaubt.'
    },
    'it': {
        'profileUpdated': 'Profilo aggiornato con successo',
        'invalidPassword': 'Password corrente non valida',
        'passwordRequired': 'Devi inserire la password corrente per cambiarla',
        'passwordMismatch': 'Le nuove password non corrispondono',
        'usernameExists': 'Il nome utente esiste già',
        'invalidFileType': 'Tipo di file non valido. Sono consentite solo immagini.'
    },
    'pt': {
        'profileUpdated': 'Perfil atualizado com sucesso',
        'invalidPassword': 'Senha atual incorreta',
        'passwordRequired': 'Você deve inserir sua senha atual para alterá-la',
        'passwordMismatch': 'As novas senhas não correspondem',
        'usernameExists': 'O nome de usuário já existe',
        'invalidFileType': 'Tipo de arquivo inválido. Apenas imagens são permitidas.'
    },
    'ru': {
        'profileUpdated': 'Профиль успешно обновлен',
        'invalidPassword': 'Неверный текущий пароль',
        'passwordRequired': 'Вы должны ввести текущий пароль для его изменения',
        'passwordMismatch': 'Новые пароли не совпадают',
        'usernameExists': 'Имя пользователя уже существует',
        'invalidFileType': 'Недопустимый тип файла. Разрешены только изображения.'
    },
    'zh': {
        'profileUpdated': '个人资料更新成功',
        'invalidPassword': '当前密码不正确',
        'passwordRequired': '您必须输入当前密码才能更改',
        'passwordMismatch': '新密码不匹配',
        'usernameExists': '用户名已存在',
        'invalidFileType': '文件类型无效。只允许图像。'
    },
    'ja': {
        'profileUpdated': 'プロフィールが正常に更新されました',
        'invalidPassword': '現在のパスワードが正しくありません',
        'passwordRequired': 'パスワードを変更するには現在のパスワードを入力する必要があります',
        'passwordMismatch': '新しいパスワードが一致しません',
        'usernameExists': 'ユーザー名はすでに存在します',
        'invalidFileType': 'ファイルタイプが無効です。画像のみ許可されています。'
    },
    'ko': {
        'profileUpdated': '프로필이 성공적으로 업데이트되었습니다',
        'invalidPassword': '현재 비밀번호가 올바르지 않습니다',
        'passwordRequired': '비밀번호를 변경하려면 현재 비밀번호를 입력해야 합니다',
        'passwordMismatch': '새 비밀번호가 일치하지 않습니다',
        'usernameExists': '사용자 이름이 이미 존재합니다',
        'invalidFileType': '잘못된 파일 형식입니다. 이미지만 허용됩니다.'
    },
    'ar': {
        'profileUpdated': 'تم تحديث الملف الشخصي بنجاح',
        'invalidPassword': 'كلمة المرور الحالية غير صحيحة',
        'passwordRequired': 'يجب عليك إدخال كلمة المرور الحالية لتغييرها',
        'passwordMismatch': 'كلمات المرور الجديدة غير متطابقة',
        'usernameExists': 'اسم المستخدم موجود بالفعل',
        'invalidFileType': 'نوع ملف غير صالح. يُسمح بالصور فقط.'
    },
    'hi': {
        'profileUpdated': 'प्रोफ़ाइल सफलतापूर्वक अपडेट किया गया',
        'invalidPassword': 'वर्तमान पासवर्ड गलत है',
        'passwordRequired': 'पासवर्ड बदलने के लिए आपको अपना वर्तमान पासवर्ड दर्ज करना होगा',
        'passwordMismatch': 'नए पासवर्ड मेल नहीं खाते',
        'usernameExists': 'उपयोगकर्ता नाम पहले से मौजूद है',
        'invalidFileType': 'अमान्य फ़ाइल प्रकार। केवल चित्र की अनुमति है।'
    },
    'nl': {
        'profileUpdated': 'Profiel succesvol bijgewerkt',
        'invalidPassword': 'Huidig wachtwoord ongeldig',
        'passwordRequired': 'U moet uw huidige wachtwoord invoeren om het te wijzigen',
        'passwordMismatch': 'Nieuwe wachtwoorden komen niet overeen',
        'usernameExists': 'Gebruikersnaam bestaat al',
        'invalidFileType': 'Ongeldig bestandstype. Alleen afbeeldingen zijn toegestaan.'
    },
    'sv': {
        'profileUpdated': 'Profil uppdaterad framgångsrikt',
        'invalidPassword': 'Ogiltigt nuvarande lösenord',
        'passwordRequired': 'Du måste anga ditt nuvarande lösenord för att ändra det',
        'passwordMismatch': 'Nya lösenord matchar inte',
        'usernameExists': 'Användarnamnet finns redan',
        'invalidFileType': 'Ogiltig filtyp. Endast bilder är tillåtna.'
    },
    'pl': {
        'profileUpdated': 'Profil zaktualizowany pomyślnie',
        'invalidPassword': 'Nieprawidłowe bieżące hasło',
        'passwordRequired': 'Musisz wprowadzić bieżące hasło, aby je zmienić',
        'passwordMismatch': 'Nowe hasła się nie zgadzają',
        'usernameExists': 'Nazwa użytkownika już istnieje',
        'invalidFileType': 'Nieprawidłowy typ pliku. Dozwolone są tylko obrazy.'
    },
    'tr': {
        'profileUpdated': 'Profil başarıyla güncellendi',
        'invalidPassword': 'Geçersiz mevcut şifre',
        'passwordRequired': 'Şifrenizi değiştirmek için mevcut şifrenizi girmelisiniz',
        'passwordMismatch': 'Yeni şifreler eşleşmiyor',
        'usernameExists': 'Kullanıcı adı zaten mevcut',
        'invalidFileType': 'Geçersiz dosya türü. Yalnızca resim dosyalarına izin verilir.'
    },
    'vi': {
        'profileUpdated': 'Cập nhật hồ sơ thành công',
        'invalidPassword': 'Mật khẩu hiện tại không chính xác',
        'passwordRequired': 'Bạn phải nhập mật khẩu hiện tại để thay đổi nó',
        'passwordMismatch': 'Mật khẩu mới không khớp',
        'usernameExists': 'Tên người dùng đã tồn tại',
        'invalidFileType': 'Loại tệp không hợp lệ. Chỉ cho phép hình ảnh.'
    },
    'th': {
        'profileUpdated': 'อัปเดตโปรไฟล์สำเร็จ',
        'invalidPassword': 'รหัสผ่านปัจจุบันไม่ถูกต้อง',
        'passwordRequired': 'คุณต้องป้อนรหัสผ่านปัจจุบันเพื่อเปลี่ยน',
        'passwordMismatch': 'รหัสผ่านใหม่ไม่ตรงกัน',
        'usernameExists': 'ชื่อผู้ใช้มีอยู่แล้ว',
        'invalidFileType': 'ประเภทไฟล์ไม่ถูกต้อง อนุญาตเฉพาะรูปภาพเท่านั้น'
    },
    'id': {
        'profileUpdated': 'Profil berhasil diperbarui',
        'invalidPassword': 'Kata sandi saat ini tidak valid',
        'passwordRequired': 'Anda harus memasukkan kata sandi saat ini untuk mengubahnya',
        'passwordMismatch': 'Kata sandi baru tidak cocok',
        'usernameExists': 'Nama pengguna sudah ada',
        'invalidFileType': 'Jenis file tidak valid. Hanya gambar yang diperbolehkan.'
    },
    'el': {
        'profileUpdated': 'Το προφίλ ενημερώθηκε επιτυχώς',
        'invalidPassword': 'Μη έγκυρος τρέχων κωδικός πρόσβασης',
        'passwordRequired': 'Πρέπει να εισαγάγετε τον τρέχοντα κωδικό πρόσβασης για να τον αλλάξετε',
        'passwordMismatch': 'Οι νέοι κωδικοί πρόσβασης δεν ταιριάζουν',
        'usernameExists': 'Το όνομα χρήστη υπάρχει ήδη',
        'invalidFileType': 'Μη έγκυρος τύπος αρχείου. Επιτρέπονται μόνο εικόνες.'
    }
}

# Actualizar todos los archivos JSON
for lang_code, translations in alert_translations.items():
    filename = f"{lang_code}.json"
    if os.path.exists(filename):
        with open(filename, 'r', encoding='utf-8') as f:
            data = json.load(f)
        
        # Agregar las traducciones de alertas
        if 'settings' not in data:
            data['settings'] = {}
        
        for key, value in translations.items():
            data['settings'][key] = value
        
        # Guardar el archivo actualizado
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        
        print(f"✓ {filename} actualizado")
    else:
        print(f"✗ {filename} no encontrado")

print("\nTraducciones de alertas agregadas a todos los idiomas!")
