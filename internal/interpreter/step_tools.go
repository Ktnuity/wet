package interpreter

import (
	"fmt"

	"github.com/ktnuity/wet/internal/tools"
	"github.com/ktnuity/wet/internal/types"
)

func (ip *Interpreter) StepTools(d *StepData) *StepResult {
	if d.token.Equals("download", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. download command failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("download command.")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. download command failed. failed to get destination value: %v", err)
		}

		vUrl, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. download command failed. failed to get url value: %v", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr(d, "failed to run step. download command failed. failed to get destination path.")
		}

		sUrl, okUrl := vUrl.String()
		if !okUrl {
			return ip.runtimeverr(d, "failed to run step. download command failed. failed to get url string.")
		}


		var result int = 1
		err = tools.ToolDownload(sUrl, pDst)
		if err != nil {
			ip.runtimevwarn("failed to use download tool: %v", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. download command failed. failure pushing value: %v", err))
		}
	} else if d.token.Equals("readfile", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. readfile command failed. stack is empty.")
		}

		ip.runtimev("readfile command.")
		vSrc, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. readfile command failed. failed to get source value: %v", err)
		}

		pSrc, okSrc := vSrc.Path()
		if !okSrc {
			return ip.runtimeverr(d, "failed to run step. readfile command failed. failed to get source path.")
		}

		data, err := tools.ToolReadfile(pSrc)
		if err != nil {
			ip.runtimevwarn("failed to use readfile tool: %v", err)
			err = ip.spush("")
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. readfile command failed. failure pushing value: %v", err))
			}

			err = ip.ipush(0)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. readfile command failed. failure pushing value: %v", err))
			}
		} else {
			err = ip.spush(data)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. readfile command failed. failure pushing value: %v", err))
			}

			err = ip.ipush(1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. readfile command failed. failure pushing value: %v", err))
			}
		}
	} else if d.token.Equals("move", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. move command failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("move command.")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. move command failed. failed to get destination value: %v", err)
		}

		vSrc, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. move command failed. failed to get source value: %v", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr(d, "failed to run step. move command failed. failed to get destination path.")
		}

		pSrc, okSrc := vSrc.Path()
		if !okSrc {
			return ip.runtimeverr(d, "failed to run step. move command failed. failed to get source path.")
		}

		var result int = 1
		err = tools.ToolMoveFile(pSrc, pDst)
		if err != nil {
			ip.runtimevwarn("failed to use move tool: %v", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. move command failed. failure pushing value: %v", err))
		}
	} else if d.token.Equals("copy", types.TokenTypeKeyword) {
		if ip.stack.Len() < 2 {
			return ip.runtimeverr(d, "failed to run step. copy command failed. stack size is %d. 2 is required.", ip.stack.Len())
		}

		ip.runtimev("copy command.")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. copy command failed. failed to get destination value: %v", err)
		}

		vSrc, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. copy command failed. failed to get source value: %v", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr(d, "failed to run step. copy command failed. failed to get destination path.")
		}

		pSrc, okSrc := vSrc.Path()
		if !okSrc {
			return ip.runtimeverr(d, "failed to run step. copy command failed. failed to get source path.")
		}

		occupied, err := tools.ToolCopyFile(pSrc, pDst)
		if err != nil {
			ip.runtimevwarn("failed to use copy tool: %v", err)
			if occupied {
				// File already exists: push false
				err = ip.ipush(0)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err))
				}
			} else {
				// Other failure: push false
				err = ip.ipush(0)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err))
				}
			}
		} else {
			// Success: push true
			err = ip.ipush(1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. copy command failed. failure pushing value: %v", err))
			}
		}
	} else if d.token.Equals("exist", types.TokenTypeKeyword) {
		if ip.stack.Len() == 0 {
			return ip.runtimeverr(d, "failed to run step. exist command failed. stack is empty.")
		}

		ip.runtimev("exist command.")
		vRes, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. exist command failed. failed to get path value: %v", err)
		}

		pRes, okRes := vRes.Path()
		if !okRes {
			return ip.runtimeverr(d, "failed to run step. exist command failed. failed to get path.")
		}

		err = tools.ToolExistFile(pRes)
		if err != nil {
			ip.runtimevwarn("failed to use exist tool: %v", err)
			err = ip.ipush(0)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. exist command failed. failure pushing value: %v", err))
			}
		} else {
			err = ip.ipush(1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. exist command failed. failure pushing value: %v", err))
			}
		}
	} else if d.token.Equals("touch", types.TokenTypeKeyword) {
		if ip.stack.Len() < 1 {
			return ip.runtimeverr(d, "failed to run step. touch command failed. stack is empty.")
		}

		ip.runtimev("touch command.")
		vPath, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. touch command failed. failed to get path value: %v", err)
		}

		pPath, okPath := vPath.Path()
		if !okPath {
			return ip.runtimeverr(d, "failed to run step. touch command failed. failed to get path.")
		}

		var result int = 1
		err = tools.ToolTouchFile(pPath)
		if err != nil {
			ip.runtimevwarn("failed to use touch tool: %v", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. touch command failed. failure pushing value: %v", err))
		}
	} else if d.token.Equals("mkdir", types.TokenTypeKeyword) {
		if ip.stack.Len() == 0 {
			return ip.runtimeverr(d, "failed to run step. mkdir command failed. stack is empty.")
		}

		ip.runtimev("mkdir command.")
		vPath, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. mkdir command failed. failed to get path value: %v", err)
		}

		pPath, okPath := vPath.Path()
		if !okPath {
			return ip.runtimeverr(d, "failed to run step. mkdir command failed. failed to get path.")
		}

		var result int = 1
		err = tools.ToolMakeDirectory(pPath)
		if err != nil {
			ip.runtimevwarn("failed to use mkdir tool: %v", err)
			result = 0
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. mkdir command failed. failure pushing value: %v", err))
		}
	} else if d.token.Equals("rm", types.TokenTypeKeyword) {
		if ip.stack.Len() == 0 {
			return ip.runtimeverr(d, "failed to run step. rm command failed. stack is empty.")
		}

		ip.runtimev("rm command.")
		vRes, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. rm command failed. failed to get path value: %v", err)
		}

		pRes, okRes := vRes.Path()
		if !okRes {
			return ip.runtimeverr(d, "failed to run step. rm command failed. failed to get path.")
		}

		err = tools.ToolRemoveFile(pRes)
		if err != nil {
			ip.runtimevwarn("failed to use rm tool: %v", err)
			err = ip.ipush(0)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. rm command failed. failure pushing value: %v", err))
			}
		} else {
			err = ip.ipush(1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. rm command failed. failure pushing value: %v", err))
			}
		}
	} else if d.token.Equals("unzip", types.TokenTypeKeyword) {
		ip.runtimev("unzip command.")
		vDst, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. unzip command failed. failed to get dst value: %v", err)
		}

		vRes, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. unzip command failed. failed to get res value: %v", err)
		}

		pDst, okDst := vDst.Path()
		if !okDst {
			return ip.runtimeverr(d, "failed to run step. unzip command failed. failed to get dst path.")
		}

		pRes, okRes := vRes.Path()
		if !okRes {
			return ip.runtimeverr(d, "failed to run step. unzip command failed. failed to get res path.")
		}

		result, err := tools.ToolUnzipFile(pDst, pRes)
		if err != nil {
			ip.runtimevwarn("failed to use unzip tool: %v", err)
			err = ip.ipush(-1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. unzip command failed. failure pushing error value: %v", err))
			}
			err = ip.ipush(-1)
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. unzip command failed. failure pushing error value: %v", err))
			}
		} else {
			err = ip.ipush(int(result.DirCount))
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. unzip command failed. failure pushing dir count: %v", err))
			}
			err = ip.ipush(int(result.FileCount))
			if err != nil {
				return StepBad(fmt.Errorf("failed to run step. unzip command failed. failure pushing file count: %v", err))
			}
		}
	} else if d.token.Equals("lsf", types.TokenTypeKeyword) {
		ip.runtimev("lsf command.")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. lsf command failed. failed to get dir value: %v", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. lsf command failed. failed to get dir path.")
		}

		var result int = 0
		count, err := tools.ToolLsf(pDir)
		if err != nil {
			ip.runtimevwarn("failed to use lsf tool: %v", err)
		} else {
			result = int(count)
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. lsf command failed. failure pushing count: %v", err))
		}
	} else if d.token.Equals("getf", types.TokenTypeKeyword) {
		ip.runtimev("getf command.")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. getf command failed. failed to get dir value: %v", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. getf command failed. failed to get dir path.")
		}

		vIdx, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. getf command failed. failed to get idx value: %v", err)
		}

		idx, ok := vIdx.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. getf command failed. failed to get idx int.")
		}

		var result string = ""
		name, err := tools.ToolGetf(idx, pDir)
		if err != nil {
			ip.runtimevwarn("failed to use getf tool: %v", err)
			ip.runtimevwarn("getf result is \"\" (empty string)")
		} else {
			result = name
		}

		err = ip.spush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. getf command failed. failure pushing name: %v", err))
		}
	} else if d.token.Equals("lsd", types.TokenTypeKeyword) {
		ip.runtimev("lsd command.")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. lsd command failed. failed to get dir value: %v", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. lsd command failed. failed to get dir path.")
		}

		var result int = 0
		count, err := tools.ToolLsd(pDir)
		if err != nil {
			ip.runtimevwarn("failed to use lsd tool: %v", err)
		} else {
			result = int(count)
		}

		err = ip.ipush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. lsd command failed. failure pushing count: %v", err))
		}
	} else if d.token.Equals("getd", types.TokenTypeKeyword) {
		ip.runtimev("getd command.")
		vDir, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. getd command failed. failed to get dir value: %v", err)
		}

		pDir, ok := vDir.Path()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. getd command failed. failed to get dir path.")
		}

		vIdx, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. getd command failed. failed to get index value: %v", err)
		}

		iIdx, ok := vIdx.Int()
		if !ok {
			return ip.runtimeverr(d, "failed to run step. getd command failed. failed to get index int.")
		}

		var result string = ""
		name, err := tools.ToolGetd(iIdx, pDir)
		if err != nil {
			ip.runtimevwarn("failed to use getd tool: %v", err)
		} else {
			result = name
		}

		err = ip.spush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. getd command failed. failure pushing result: %v", err))
		}
	} else if d.token.Equals("concat", types.TokenTypeKeyword) {
		ip.runtimev("concat command.")
		vB, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. concat command failed. failed to get second value: %v", err)
		}

		vA, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. concat command failed. failed to get first value: %v", err)
		}

		if pA, okA := vA.Path(); okA {
			if sB, okB := vB.String(); okB {
				result := tools.ToolConcatPath(pA, sB)
				err = ip.ppush(result)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. concat command failed. failure pushing path result: %v", err))
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. concat command failed. path concat requires string as second value.")
			}
		} else if sA, okA := vA.String(); okA {
			if sB, okB := vB.String(); okB {
				result := tools.ToolConcatString(sA, sB)
				err = ip.spush(result)
				if err != nil {
					return StepBad(fmt.Errorf("failed to run step. concat command failed. failure pushing string result: %v", err))
				}
			} else {
				return ip.runtimeverr(d, "failed to run step. concat command failed. string concat requires string as second value.")
			}
		} else {
			return ip.runtimeverr(d, "failed to run step. concat command failed. first value must be path or string.")
		}
	} else if d.token.Equals("tostring", types.TokenTypeKeyword) {
		ip.runtimev("tostring command.")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. tostring command failed. failed to get value: %v", err)
		}

		var result string
		if i, ok := v.Int(); ok {
			result = tools.ToolToStringInt(i)
		} else if s, ok := v.String(); ok {
			result = tools.ToolToStringString(s)
		} else if p, ok := v.Path(); ok {
			result = tools.ToolToStringPath(p)
		} else {
			result = ""
		}

		err = ip.spush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. tostring command failed. failure pushing string result: %v", err))
		}
	} else if d.token.Equals("token", types.TokenTypeKeyword) {
		ip.runtimev("token command.")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. d.token command failed. failed to get value: %v", err)
		}

		var result string
		if s, ok := v.String(); ok {
			result = tools.ToolToToken(s)
		} else {
			result = ""
		}

		err = ip.ppush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. d.token command failed. failure pushing path result: %v", err))
		}
	} else if d.token.Equals("absolute", types.TokenTypeKeyword) {
		ip.runtimev("absolute command.")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. absolute command failed. failed to get value: %v", err)
		}

		var result string
		if s, ok := v.String(); ok {
			result = tools.ToolToAbsolute(s)
		} else {
			result = ""
		}

		err = ip.ppush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. absolute command failed. failure pushing path result: %v", err))
		}
	} else if d.token.Equals("relative", types.TokenTypeKeyword) {
		ip.runtimev("relative command.")
		v, err := ip.pop()
		if err != nil {
			return ip.runtimeverr(d, "failed to run step. relative command failed. failed to get value: %v", err)
		}

		var result string
		if s, ok := v.String(); ok {
			result = tools.ToolToRelative(s)
		} else {
			result = ""
		}

		err = ip.ppush(result)
		if err != nil {
			return StepBad(fmt.Errorf("failed to run step. relative command failed. failure pushing path result: %v", err))
		}
	} else {
		return nil
	}

	return StepOk()
}
