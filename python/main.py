# 定义 token 类型
class Token:
    def __init__(self, token_type, value,pos):
        self.token_type = token_type
        self.value = value
        self.pos=pos

# 定义 pos 类型 也就是 存储位置信息
class Pos:
    def __init__(self,offset,line,col) -> None:
        self.offset=offset
        self.line=line
        self.col=col

# 定义词法分析器类
class Scanner:
    def __init__(self, input_string):
        self.input_string = input_string
        self.position = 0
        self.line = 1
        self.col = 1
        self.keywords ={'break','case','chan','const','continue','default','defer','else','fallthrough','for','func','go','goto','if','import','interface','map','package','range','return','select','struct','switch','type','var'}

    # 词法分析函数
    def scan(self):
        tokens = []
        while self.position < len(self.input_string):
            char = self.input_string[self.position]
            if char.isspace():
                self.update_pos()
            elif char.isalpha() or char == '_':
                token = self.extract_identifier()
                if token.value in self.keywords:
                    token.token_type = 'KEYWORD'
                tokens.append(token)
            elif char.isdigit():
                tokens.append(self.extract_number())
            elif char == '"' or char=='`':
                tokens.append(self.extract_string_and_byte())
            elif char == "'":
                tokens.append(self.extract_string_and_byte())
            elif char in '+-*%=()&|^<>!.:;/{}|[]^<>\\,':# 某些字符可能没有嵌入 后续完善
                tokens.append(self.extract_operator())
            else:
                raise ValueError(f"错误/未知字符 '{char}' 出现在 {self.position}")
        return tokens
    
    # 词法分析函数
    def scan_format(self):
        # token=None
        # 输出标题
        print("Position".ljust(15),"Type".ljust(15), "Value")
        while self.position < len(self.input_string):
            char = self.input_string[self.position]
            if char.isspace():
                self.update_pos()
                continue
            elif char.isalpha() or char == '_':
                token = self.extract_identifier()
                if token.value in self.keywords:
                    token.token_type = 'KEYWORD'
            elif char.isdigit():
                token=self.extract_number()
            elif char == '"' or char=='`' or char == "'":
                token=self.extract_string_and_byte()
            elif char in '+-*%=()&|^<>!.:;/{}|[]^<>\\,':# 某些字符可能没有嵌入 后续完善
                token=self.extract_operator()
            else:
                raise ValueError(f"错误/未知字符 '{char}' 出现在 {self.position}")
            token_type = str(token.token_type)
            value = str(token.value)
            position = f"{token.pos.line}:{token.pos.col}"
            print(position.ljust(15),token_type.ljust(15), value.ljust(15),)
    # 更新位置信息函数
    def update_pos(self):
        if self.input_string[self.position]== '\n':
            self.line += 1
            self.col = 1
        else:
            self.col += 1
        self.position+=1

    # 提取标识符
    def extract_identifier(self):
        start_position = self.position
        start_line=self.line    
        start_col=self.col
        while self.position < len(self.input_string) and (self.input_string[self.position].isalnum() or self.input_string[self.position] == '_'):
            self.update_pos()
        value = self.input_string[start_position:self.position]
        pos = Pos(start_position, start_line,start_col)
        
        return Token('KEYWORD' if value in self.keywords else 'IDENT', value,pos)
    
    # 对于操作符号 进行提取
    def extract_operator(self):
        start_position = self.position
        start_line=self.line    
        start_col=self.col
        operator = self.input_string[self.position]
        self.update_pos()
        if operator == '+':
            # ++
            if self.input_string[self.position] == '+':
                operator += '+'
                self.update_pos()
            # +=
            elif self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '-':
            # --
            if self.input_string[self.position] == '-':
                operator += '-'
                self.update_pos()
            # -=
            elif self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        
        elif operator == '*':
            # *=
            if self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '/':
            # // 单行注释
            if self.input_string[self.position] == '/':
                return self.extract_single_line_comment()
            # /* 多行注释
            elif self.input_string[self.position] == '*':
                return self.extract_multi_line_comment()
            # /= 
            elif self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '%':
            # %=
            if self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '&':
            # &^
            if self.input_string[self.position] == '^':
                operator += '^'
                self.update_pos()
            # &=
            elif self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '|':
            # |=
            if self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '^':
            # ^=
            if self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '<':
            # <<
            if self.input_string[self.position] == '<':
                operator += '<'
                self.update_pos()
                # <<=
                if self.input_string[self.position] == '=':
                    operator += '='
                    self.update_pos()
            # <=
            elif self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        elif operator == '>':
            # >>
            if self.input_string[self.position] == '>':
                operator += '>'
                self.update_pos()
                # >>=
                if self.input_string[self.position] == '=':
                    operator += '='
                    self.update_pos()
            # >=
            elif self.input_string[self.position] == '=':
                operator += '='
                self.update_pos()
        value = operator
        pos = Pos(start_position, start_line,start_col)
        return Token('OP/SEP', value,pos)
    
    # 提取数字 分为float/int/imag 补充实现提取虚数的情况
    def extract_number(self):
        start_position = self.position
        start_line = self.line
        start_col = self.col
        
        # 定义进制标识符和对应的进制
        prefixes = {'0x': 16, '0o': 8, '0b': 2}
        # 检查正负号
        if self.input_string[self.position]=='-':
            self.update_pos()
        # 检查是否存在进制标识符
        for prefix, base in prefixes.items():
            if self.input_string[self.position:self.position+len(prefix)] == prefix:
                self.position += len(prefix)  # 跳过进制标识符
                break
        
        # 提取整数部分
        while self.position < len(self.input_string) and self.input_string[self.position].isdigit():
            self.update_pos()
        
        # 提取小数部分
        if self.position < len(self.input_string) and self.input_string[self.position] == '.':
            self.update_pos()  # 跳过小数点
            while self.position < len(self.input_string) and self.input_string[self.position].isdigit():
                self.update_pos()
        
        # 提取科学计数法标识的数字
        if self.position < len(self.input_string) and (self.input_string[self.position] == 'e' or self.input_string[self.position] == 'E'):
            self.update_pos()  # 跳过'e'或'E'
            if self.position < len(self.input_string) and (self.input_string[self.position] == '+' or self.input_string[self.position] == '-'):
                self.update_pos()  # 跳过正负号
            while self.position < len(self.input_string) and self.input_string[self.position].isdigit():
                self.update_pos()
        
        # 判断是不是 虚数 注意 虚数识别仅仅需要识别到他的复数部分就行了 
        is_imag=False
        if self.position < len(self.input_string) and (self.input_string[self.position] == 'i'):
            is_imag=True
        is_float=False
        # 计算数字
        pos = Pos(start_position, start_line, start_col)
        num_str = self.input_string[start_position:self.position]
        if '.' in num_str or 'e' in num_str or 'E' in num_str:  # 浮点数或科学计数法标识的数字
            num = float(num_str)
            is_float=True
        elif 'x' in num_str or 'X' in num_str:  # 十六进制
            num = int(num_str, 16)
        elif 'o' in num_str or 'O' in num_str:  # 八进制
            num = int(num_str, 8)
        elif 'b' in num_str or 'B' in num_str:  # 二进制
            num = int(num_str, 2)
        else:  # 十进制整数
            num = int(num_str)
        if is_imag:
            self.update_pos()
            return Token('IMAG',self.input_string[start_position:self.position],pos)
        if is_float:
            return Token('FLOAT', num, pos)
        return Token('INT', num, pos)

    def extract_string_and_byte(self):
        pre = self.input_string[self.position]
        start_position = self.position 
        start_line = self.line    
        start_col = self.col
        self.update_pos()
        escape = False
        value =pre
        while self.position < len(self.input_string) and (self.input_string[self.position] != pre or escape):
            if escape:
                value += '\\'+self.input_string[self.position]
                escape = False
            else:
                if self.input_string[self.position] == '\\':
                    escape = True
                else:
                    value +=self.input_string[self.position]
            self.update_pos()

        if self.position == len(self.input_string):
            raise ValueError("字符串解析存在错误")
            
        self.update_pos()
        pos = Pos(start_position, start_line, start_col)
        value = value.replace("\n",'\\n').strip()# 对于 \n进行转义处理
        return Token("BYTE" if pre == "'" else 'STRING', value+pre, pos)

    # 提取单行注释
    def extract_single_line_comment(self):
        start_position = self.position-1
        start_line=self.line    
        start_col=self.col-1
        
        while self.position < len(self.input_string) and self.input_string[self.position] != '\n':
            self.update_pos()
        value = self.input_string[start_position:self.position].strip()
        pos = Pos(start_position, start_line,start_col)
        return Token('COMMENT', value,pos)
    
    # 提取多行注释
    def extract_multi_line_comment(self):
        start_position = self.position-1
        start_line=self.line
        start_col=self.col-1
        while self.position < len(self.input_string) - 1 and self.input_string[self.position:self.position + 2] != '*/':
            self.update_pos()
        if self.position == len(self.input_string) - 1:
            raise ValueError("多行注释出现错误")
        
        # 对于字符中存在的换行进行转义处理 避免换行显示出来 
        pos = Pos(start_position, start_line,start_col)
        self.update_pos()
        self.update_pos()
        value = self.input_string[start_position:self.position].replace("\n",'\\n').replace("\r",'\\r').strip()
        return Token('COMMENT', value,pos)
    
def test_lexer(input_string):
    scanner = Scanner(input_string)
    tokens = scanner.scan()

    # 输出标题
    print("Position".ljust(15),"Type".ljust(15), "Value")

    # 输出词法分析结果
    for token in tokens:
        token_type = str(token.token_type)
        value = str(token.value)
        position = f"{token.pos.line}:{token.pos.col}"
        print(position.ljust(15),token_type.ljust(15), value.ljust(15),)
def test_lexer_format(input_string):
    scanner = Scanner(input_string)
    scanner.scan_format()
# 测试
if __name__ == "__main__":
    input_string =open("./go/main.go",encoding="utf-8").read()
    # test_lexer(input_string)
    test_lexer_format(input_string)
