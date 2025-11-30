package argsync

func GenObjectArgTypeCode(name string, index string) (string, string) {
	code := `
  Object wg__` + name + ` = Object::New(wg_env);
  if (wg_info.Length() > ` + index + `) {
    wg__` + name + ` = wg_info[` + index + `].As<Object>();
  }
  string wg_` + name + ` = wg_object_to_string(wg__` + name + `, wg_env);
  unique_ptr<char[]> wg_` + name + `_buf(new char[wg_` + name + `.length() + 1]);
  strcpy(wg_` + name + `_buf.get(), wg_` + name + `.c_str());
  char *` + name + ` = wg_` + name + `_buf.get();`

	endCode := ""
	return code, endCode
}
