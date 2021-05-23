package helper
import(
	"fmt"
	"unicode"
)
func CheckPassword(password string)(isValid bool,errors []string){
	isValid=false
	errors =[]string{}
	hasLowerCase :=false
	hasUpperCase :=false
	hasSpecialChar :=false
	isLongEnough :=false
	specialChars:=[]rune{'!','#','$','%','&','*','+','@','='}
	fmt.Println(specialChars)
	length:=len(password)
	for i:=0;i<length;i++{
		char:=rune(password[i])
		if unicode.IsLower(char) && char != ' ' {
			hasLowerCase=true;
		}
		if unicode.IsUpper(char) && char != ' ' {
			hasUpperCase=true;
		}
		for j:=0;j<len(specialChars);j++{
			specialChar:=rune(specialChars[j])
			if char==specialChar{
				hasSpecialChar=true;
			}
		}
	}
	if length>=8 && length<=16{
		isLongEnough =true
	}
	if hasLowerCase && hasUpperCase && hasSpecialChar && isLongEnough {
		isValid=true
	}else{
		if !hasLowerCase{
			errors=append(errors,"must have at least one lower case letter")
		}
		if !hasUpperCase{
			errors=append(errors,"must have at least one upper case letter")
		}
		if !hasSpecialChar{
			errors=append(errors,"must have at least one special Character")
		}
		if !isLongEnough{
			errors=append(errors,"must be between 8 and 16 characters")
		}
	}
	return
}