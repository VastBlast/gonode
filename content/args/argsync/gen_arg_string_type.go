package argsync

func GenStringArgTypeCode(name string, index string) (string, string) {
	code := `
  string wg_` + name + ` = "";
  if (wg_info.Length() > ` + index + `) {
    wg_` + name + ` = wg_info[` + index + `].As<String>().Utf8Value();
  }
  unique_ptr<char[]> wg_` + name + `_buf(new char[wg_` + name + `.length() + 1]);
  strcpy(wg_` + name + `_buf.get(), wg_` + name + `.c_str());
  char *` + name + ` = wg_` + name + `_buf.get();`

	endCode := ""
	return code, endCode
}
