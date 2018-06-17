package gofpdf

import (
	"log"
	"strconv"

	"golang.org/x/crypto/pkcs12"
)

type SignatureInfo struct {
	Name        string
	Location    string
	Reason      string
	ContactInfo string
}

// SetSignature Enable document signature.
// The digital signature improve document authenticity and integrity and allows o enable extra features on Acrobat Reader.
// To create self-signed signature: openssl req -x509 -nodes -days 365000 -newkey rsa:1024 -keyout tcpdf.crt -out tcpdf.crt
// To export crt to p12: openssl pkcs12 -export -in tcpdf.crt -out tcpdf.p12
// To convert pfx certificate to pem: openssl pkcs12 -in tcpdf.pfx -out tcpdf.crt -nodes
// $extracerts (string) specifies the name of a file containing a bunch of extra certificates to include in the signature which can for example be used to help the recipient to verify the certificate that you used.
// CertType The access permissions granted for this document. Valid values shall be: 1 = No changes to the document shall be permitted; any change to the document shall invalidate the signature; 2 = Permitted changes shall be filling in forms, instantiating page templates, and signing; other changes shall invalidate the signature; 3 = Permitted changes shall be the same as for 2, as well as annotation creation, deletion, and modification; other changes shall invalidate the signature.
// Info array of option information: Name, Location, Reason, ContactInfo.
// Approval Enable approval signature eg. for PDF incremental update
func (f *Fpdf) SetSignature(Pfx []byte, PfxPassword string, CertType int, Info SignatureInfo, Approval string) {
	priv, cert, err := pkcs12.Decode(Pfx, PfxPassword)
	if err != nil {
		log.Fatal(err)
	}
	f.signature = true
	f.n = f.n + 1
	f.signObjID = f.n // signature widget
	f.n = f.n + 1     // signature object
	// if (strlen($signing_cert) == 0) {
	// 	$this->Error('Please provide a certificate file and password!');
	// }
	// if (strlen($private_key) == 0) {
	// 	$private_key = $signing_cert;
	// }
	signatureData := make(map[string]interface{})
	//signatureData["signcert"] = $signing_cert
	//signatureData['privkey'] = $private_key
	//signatureData['password'] = $private_key_password
	//signatureData['extracerts'] = $extracerts
	//signatureData['cert_type'] = $cert_type
	signatureData["info"] = Info
	signatureData["approval"] = Approval
	f.signatureData = signatureData

}

/**
 * Get the array that defines the signature appearance (page and rectangle coordinates).
 * $x (float) Abscissa of the upper-left corner.
 * $y (float) Ordinate of the upper-left corner.
 * $w (float) Width of the signature area.
 * $h (float) Height of the signature area.
 * $page (int) option page number (if < 0 the current page is used).
 * $name (string) Name of the signature.
 * @return (array) Array defining page and rectangle coordinates of signature appearance.
 */
func (f *Fpdf) getSignatureAppearance(x, y, w, h float64, page int, name string) *signAppearance {
	sigapp := new(signAppearance)
	totalpages, _ := strconv.Atoi(f.aliasNbPagesStr)
	if (page < 1) || (page > totalpages) {
		sigapp.page = f.page
	} else {
		sigapp.page = page
	}
	if name == "" {
		sigapp.name = "Signature"
	} else {
		sigapp.name = name
	}
	//a := x * f.k;
	//b := f.pagedim[(sigapp.page)]['h'] - ((y + h) * f.k);
	//c := w * f.k;
	//d := h * f.k;
	//sigapp.rect = sprintf("%F %F %F %F", a, b, (a + c), (b + d))
	return sigapp
}

// SetTimeStamp Enable document timestamping .
// The trusted timestamping improve document security that means that no one should be able to change the document once it has been recorded.
// Use with digital signature only!
// $tsaHost (string) Time Stamping Authority (TSA) server (prefixed with 'https://')
// $tsaUsername (string) Specifies the username for TSA authorization (optional) OR specifies the TSA authorization PEM file (see: example_66.php, optional)
// $tsaPassword (string) Specifies the password for TSA authorization (optional)
// tsaCert (string) Specifies the location of TSA certificate for authorization (optional for cURL)
func (f *Fpdf) SetTimeStamp(tsaHost, tsaUsername, tsaPassword, tsaCert string) {
	// $this->tsa_data = array();
	// if (!function_exists('curl_init')) {
	// 	$this->Error('Please enable cURL PHP extension!');
	// }
	// if (strlen($tsa_host) == 0) {
	// 	$this->Error('Please specify the host of Time Stamping Authority (TSA)!');
	// }
	// $this->tsa_data['tsa_host'] = $tsa_host;
	// if (is_file($tsa_username)) {
	// 	$this->tsa_data['tsa_auth'] = $tsa_username;
	// } else {
	// 	$this->tsa_data['tsa_username'] = $tsa_username;
	// }
	// $this->tsa_data['tsa_password'] = $tsa_password;
	// $this->tsa_data['tsa_cert'] = $tsa_cert;
	// $this->tsa_timestamp = true;
}

// NOT YET IMPLEMENTED
// Request TSA for a timestamp
// signature (string) Digital signature as binary string
// @return (string) Timestamped digital signature
func (f *Fpdf) applyTSA(signature string) string {
	if !f.tsaTimestamp {
		return signature
	}
	//@TODO: implement this feature
	return signature
}

// SetSignatureAppearance set the digital signature appearance (a cliccable rectangle area to get signature properties)
// x (float) Abscissa of the upper-left corner.
// y (float) Ordinate of the upper-left corner.
// w (float) Width of the signature area.
// h (float) Height of the signature area.
// page (int) option page number (if < 0 the current page is used).
// name (string) Name of the signature.
func (f *Fpdf) SetSignatureAppearance(x, y, w, h float64, page int, name string) {
	f.signatureAppearance = f.getSignatureAppearance(x, y, w, h, page, name)
}

// AddEmptySignatureAppearance Add an empty digital signature appearance (a cliccable rectangle area to get signature properties)
// x (float) Abscissa of the upper-left corner.
// y (float) Ordinate of the upper-left corner.
// w (float) Width of the signature area.
// h (float) Height of the signature area.
// page (int) option page number (if < 0 the current page is used).
// name (string) Name of the signature.
func (f *Fpdf) AddEmptySignatureAppearance(x, y, w, h float64, page int, name string) {
	f.n = f.n + 1
	//f.emptySignatureAppearance = array('objid' => f.n) + f.getSignatureAppearanceArray(x, y, w, h, page, name)
}

// Add certification signature (DocMDP or UR3)
// You can set only one signature type
func (f *Fpdf) putsignature() {
	// if ((!f.sign) || (!isset($this->signature_data['cert_type']))) {
	// 	return;
	// }
	// $sigobjid = ($this->sig_obj_id + 1);
	// $out = $this->_getobj($sigobjid)."\n";
	// $out .= "<< /Type /Sig';
	// $out .= " /Filter /Adobe.PPKLite';
	// $out .= " /SubFilter /adbe.pkcs7.detached';
	// $out .= " ".TCPDF_STATIC::$byterange_string;
	// $out .= " /Contents<".str_repeat('0', $this->signature_max_length).">";
	// if (empty($this->signature_data['approval']) OR ($this->signature_data['approval'] != 'A')) {
	// 	$out .= ' /Reference ['; // array of signature reference dictionaries
	// 	$out .= ' << /Type /SigRef';
	// 	if ($this->signature_data['cert_type'] > 0) {
	// 		$out .= ' /TransformMethod /DocMDP';
	// 		$out .= ' /TransformParams <<';
	// 		$out .= ' /Type /TransformParams';
	// 		$out .= ' /P '.$this->signature_data['cert_type'];
	// 		$out .= ' /V /1.2';
	// 	} else {
	// 		$out .= ' /TransformMethod /UR3';
	// 		$out .= ' /TransformParams <<';
	// 		$out .= ' /Type /TransformParams';
	// 		$out .= ' /V /2.2';
	// 		if (!TCPDF_STATIC::empty_string($this->ur['document'])) {
	// 			$out .= ' /Document['.$this->ur['document'].']';
	// 		}
	// 		if (!TCPDF_STATIC::empty_string($this->ur['form'])) {
	// 			$out .= ' /Form['.$this->ur['form'].']';
	// 		}
	// 		if (!TCPDF_STATIC::empty_string($this->ur['signature'])) {
	// 			$out .= ' /Signature['.$this->ur['signature'].']';
	// 		}
	// 		if (!TCPDF_STATIC::empty_string($this->ur['annots'])) {
	// 			$out .= ' /Annots['.$this->ur['annots'].']';
	// 		}
	// 		if (!TCPDF_STATIC::empty_string($this->ur['ef'])) {
	// 			$out .= ' /EF['.$this->ur['ef'].']';
	// 		}
	// 		if (!TCPDF_STATIC::empty_string($this->ur['formex'])) {
	// 			$out .= ' /FormEX['.$this->ur['formex'].']';
	// 		}
	// 	}
	// 	$out .= ' >>'; // close TransformParams
	// 	// optional digest data (values must be calculated and replaced later)
	// 	//$out .= ' /Data ********** 0 R';
	// 	//$out .= ' /DigestMethod/MD5';
	// 	//$out .= ' /DigestLocation[********** 34]';
	// 	//$out .= ' /DigestValue<********************************>';
	// 	out += " >>";
	// 	out += " ]"; // end of reference
	// }
	// if (isset($this->signature_data['info']['Name']) AND !TCPDF_STATIC::empty_string($this->signature_data['info']['Name'])) {
	// 	$out .= ' /Name '.$this->_textstring($this->signature_data['info']['Name'], $sigobjid);
	// }
	// if (isset($this->signature_data['info']['Location']) AND !TCPDF_STATIC::empty_string($this->signature_data['info']['Location'])) {
	// 	$out .= ' /Location '.$this->_textstring($this->signature_data['info']['Location'], $sigobjid);
	// }
	// if (isset($this->signature_data['info']['Reason']) AND !TCPDF_STATIC::empty_string($this->signature_data['info']['Reason'])) {
	// 	$out .= ' /Reason '.$this->_textstring($this->signature_data['info']['Reason'], $sigobjid);
	// }
	// if (isset($this->signature_data['info']['ContactInfo']) AND !TCPDF_STATIC::empty_string($this->signature_data['info']['ContactInfo'])) {
	// 	$out .= ' /ContactInfo '.$this->_textstring($this->signature_data['info']['ContactInfo'], $sigobjid);
	// }
	// $out .= ' /M '.$this->_datestring($sigobjid, $this->doc_modification_timestamp);
	// $out .= ' >>';
	// $out .= "\n".'endobj';
	// $this->_out($out);
}
