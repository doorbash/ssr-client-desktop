import React, { Component, useContext } from "react"
import { useLocation, useNavigate } from "react-router-dom"
import "tw-elements"
import AppContext from "../../contexts/AppContext"
import "./NewProxy.css"

class NewProxy extends Component {
  constructor(props) {
    super(props)
    if (props.location.pathname === "/edit-proxy") {
      props.setConfig("Edit Proxy", true, "white")
    } else {
      props.setConfig("New Proxy", true, "white")
    }
    this.handleSubmit = this.handleSubmit.bind(this)

    if (props.location.state) {
      this.state = {
        ...props.location.state,
        error: undefined,
      }
    } else {
      this.state = {
        name: "My Proxy",
        s: "",
        p: 8388,
        b: "127.0.0.1",
        l: 1080,
        k: "",
        m: "aes-256-cfb",
        o: "http_simple",
        op: "",
        oo: "origin",
        oop: "",
        t: 10,
        f: "",
        error: undefined,
      }
    }
  }

  handleSubmit(e) {
    e.preventDefault()
    console.log("form submitted!")
    this.setState({
      error: undefined,
    })
    if (!window.go) return
    ;(async () => {
      try {
        if (this.props.location.pathname === "/edit-proxy") {
          await window.go.main.App.UpdateProxy({
            ...this.state,
            create_time: Date.now(),
          })
          this.props.navigate(-1)
        } else {
          await window.go.main.App.InsertProxy({
            ...this.state,
            create_time: Date.now(),
            status: 1,
          })
          this.props.navigate(-1)
        }
      } catch (e) {
        this.setState({
          error: e,
        })
      }
    })()
  }

  render() {
    return (
      <>
        <form
          className="px-8 pt-4 mb-4 grid grid-cols-2 items-center gap-3 place-items-start"
          onSubmit={this.handleSubmit}
          autoComplete="off"
        >
          <label className="new-proxy-label" htmlFor="proxy-name">
            Proxy name:
          </label>
          <input
            className="new-proxy-input"
            id="proxy-name"
            required={true}
            defaultValue={this.state.name}
            onChange={(e) => {
              this.setState({ name: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-s">
            Server address:
          </label>
          <input
            className="new-proxy-input"
            id="-s"
            required={true}
            defaultValue={this.state.s}
            onChange={(e) => {
              this.setState({ s: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-p">
            Server port:
          </label>
          <input
            type="number"
            className="new-proxy-input"
            id="-p"
            required={true}
            defaultValue={this.state.p}
            onChange={(e) => {
              this.setState({ p: e.target.valueAsNumber })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-b">
            Local binding address:
          </label>
          <input
            className="new-proxy-input"
            id="-b"
            required={true}
            defaultValue={this.state.b}
            onChange={(e) => {
              this.setState({ b: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-l">
            Local port:
          </label>
          <input
            type="number"
            className="new-proxy-input"
            id="-l"
            required={true}
            defaultValue={this.state.l}
            onChange={(e) => {
              this.setState({ l: e.target.valueAsNumber })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-k">
            Password:
          </label>
          <input
            type="password"
            className="new-proxy-input"
            id="-k"
            autoComplete="new-password"
            required={true}
            defaultValue={this.state.k}
            onChange={(e) => {
              this.setState({ k: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-m">
            Encryption method:
          </label>
          <input
            className="new-proxy-input"
            id="-m"
            required={true}
            defaultValue={this.state.m}
            onChange={(e) => {
              this.setState({ m: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-o">
            Obfsplugin:
          </label>
          <input
            className="new-proxy-input"
            id="-o"
            required={true}
            defaultValue={this.state.o}
            onChange={(e) => {
              this.setState({ o: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="--op">
            Obfs param:
          </label>
          <input
            className="new-proxy-input"
            id="--op"
            defaultValue={this.state.op}
            onChange={(e) => {
              this.setState({ op: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-O">
            Protocol:
          </label>
          <input
            className="new-proxy-input"
            id="-O"
            required={true}
            defaultValue={this.state.oo}
            onChange={(e) => {
              this.setState({ oo: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="--Op">
            Protocol param:
          </label>
          <input
            className="new-proxy-input"
            id="--Op"
            defaultValue={this.state.oop}
            onChange={(e) => {
              this.setState({ oop: e.target.value })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-t">
            Socket timeout (seconds):
          </label>
          <input
            type="number"
            className="new-proxy-input"
            id="-t"
            required={true}
            defaultValue={this.state.t}
            onChange={(e) => {
              this.setState({ t: e.target.valueAsNumber })
            }}
          ></input>

          <label className="new-proxy-label" htmlFor="-f">
            Socks5 forward proxy address: <br />{" "}
            <span>(example: 127.0.0.1:8080)</span>
          </label>
          <input
            className="new-proxy-input"
            id="-f"
            defaultValue={this.state.f}
            onChange={(e) => {
              this.setState({ f: e.target.value })
            }}
          />

          <div className="new-proxy-submit">
            <input className="btn-primary" type="submit" value="Save" />
          </div>
        </form>
        {this.state.error && (
          <div
            className="error-alert alert bg-red-100 rounded-lg py-5 px-6 mb-1 text-base text-red-700 flex items-center w-full fixed right-0 left-0 bottom-0"
            role="alert"
          >
            <strong className="mr-1">Error</strong> while inserting to database:{" "}
            {this.state.error}
          </div>
        )}
      </>
    )
  }
}

export default (props) => {
  return (
    <NewProxy
      {...props}
      navigate={useNavigate()}
      location={useLocation()}
      {...useContext(AppContext)}
    />
  )
}
