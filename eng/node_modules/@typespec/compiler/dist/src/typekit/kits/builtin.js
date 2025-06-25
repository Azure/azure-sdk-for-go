import { defineKit } from "../define-kit.js";
defineKit({
    builtin: {
        get string() {
            return this.program.checker.getStdType("string");
        },
        get boolean() {
            return this.program.checker.getStdType("boolean");
        },
        get bytes() {
            return this.program.checker.getStdType("bytes");
        },
        get decimal() {
            return this.program.checker.getStdType("decimal");
        },
        get decimal128() {
            return this.program.checker.getStdType("decimal128");
        },
        get duration() {
            return this.program.checker.getStdType("duration");
        },
        get float() {
            return this.program.checker.getStdType("float");
        },
        get float32() {
            return this.program.checker.getStdType("float32");
        },
        get float64() {
            return this.program.checker.getStdType("float64");
        },
        get int8() {
            return this.program.checker.getStdType("int8");
        },
        get int16() {
            return this.program.checker.getStdType("int16");
        },
        get int32() {
            return this.program.checker.getStdType("int32");
        },
        get int64() {
            return this.program.checker.getStdType("int64");
        },
        get integer() {
            return this.program.checker.getStdType("integer");
        },
        get offsetDateTime() {
            return this.program.checker.getStdType("offsetDateTime");
        },
        get plainDate() {
            return this.program.checker.getStdType("plainDate");
        },
        get plainTime() {
            return this.program.checker.getStdType("plainTime");
        },
        get safeInt() {
            return this.program.checker.getStdType("safeint");
        },
        get uint8() {
            return this.program.checker.getStdType("uint8");
        },
        get uint16() {
            return this.program.checker.getStdType("uint16");
        },
        get uint32() {
            return this.program.checker.getStdType("uint32");
        },
        get uint64() {
            return this.program.checker.getStdType("uint64");
        },
        get url() {
            return this.program.checker.getStdType("url");
        },
        get utcDateTime() {
            return this.program.checker.getStdType("utcDateTime");
        },
        get numeric() {
            return this.program.checker.getStdType("numeric");
        },
    },
});
//# sourceMappingURL=builtin.js.map