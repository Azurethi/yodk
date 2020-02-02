import * as vscode from 'vscode'
import * as assert from 'assert'
import {getExePath, getContext, runYodkCommand} from '../extension'
import { activate, getDocUri } from './helper'

describe('Interact with binary', async () => {

  it('Find on linux', async () => {
    await activate(getDocUri("correct.yolol"))
    let path = getExePath("linux")
    assert.equal(path,getContext().asAbsolutePath("./bin/yodk"))
  })

  it('Find on Windows', async () => {
    await activate(getDocUri("correct.yolol"))
    let path = getExePath("win32")
    assert.equal(path,getContext().asAbsolutePath("./bin/yodk.exe"))
  })

  it('Find with env var', async () => {
    process.env["YODK_EXECUTABLE"] = "/test/path/yodk"
    await activate(getDocUri("correct.yolol"))
    let path = getExePath("win32")
    process.env["YODK_EXECUTABLE"] = ""
    assert.equal(path,"/test/path/yodk")
  })

  it('Answers on version', async () =>{
    await activate(getDocUri("correct.yolol"))
    let result = await runYodkCommand(["version"])
    let correct = result["output"] == "\nUNVERSIONED BUILD\n" || result["output"].startsWith("\nv")
    assert.equal(result["code"],0)
    assert.equal(correct,true)
  })

  it('Errors on unknown', async () =>{
    await activate(getDocUri("correct.yolol"))
    let result = await runYodkCommand(["unknown","command"])
    assert.equal(result["code"],1)
  })

})