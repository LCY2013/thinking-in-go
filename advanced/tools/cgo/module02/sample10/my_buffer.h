#include <string>

struct MyBuffer {
    std::string* s_;

    MyBuffer(int size) {
        this->s_ = new std::string(size, char('\0'));
    }

    ~MyBuffer() {
        delete this->s_;
    }

    int size() const {
        return this->s_->size();
    }

    char* Data() {
        return (char*)this->s_->data();
    }
}

int main() {
    auto pBuf = new MyBuffer(1024);

    auto data = pBuf->Data();
    auto size = pBuf->size();

    delete pBuf;
}
