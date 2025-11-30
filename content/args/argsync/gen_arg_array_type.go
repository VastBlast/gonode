package argsync

func GenArrayArgTypeCode(name string, index string) (string, string) {
	code := `
  Array wg__` + name + ` = Array::New(wg_env);
  if (wg_info.Length() > ` + index + `) {
    wg__` + name + ` = wg_info[` + index + `].As<Array>();
  }
  string wg_` + name + ` = wg_array_to_string(wg__` + name + `, wg_env);
  unique_ptr<char[]> wg_` + name + `_buf(new char[wg_` + name + `.length() + 1]);
  strcpy(wg_` + name + `_buf.get(), wg_` + name + `.c_str());
  char *` + name + ` = wg_` + name + `_buf.get();`

	endCode := ""
	return code, endCode
}
