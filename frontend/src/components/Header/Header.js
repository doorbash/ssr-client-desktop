import { Settings, GitHub, ArrowBack } from "@material-ui/icons"
import { useLocation, useNavigate } from "react-router-dom"
import "./Header.css"

import React, { Component, useContext } from "react"
import AppContext from "../../contexts/AppContext"

class Header extends Component {
  componentDidUpdate() {
    this.props.updateTooltips()
  }

  render() {
    return (
      <header className="header">
        <div className="flex justify-start items-center gap-5">
          {this.props.showBackButton && (
            <ArrowBack
              fontSize="large"
              className="menu-icon"
              onClick={() => {
                this.props.navigate(-1)
              }}
            />
          )}
          <div className="text-[1.35rem] truncate max-w-[10em]">
            {this.props.title}
          </div>
        </div>
        <div className="flex justify-end items-center gap-3">
          <div
            data-bs-toggle="tooltip"
            title="SETTINGS"
            data-bs-placement="bottom"
            className="menu-icon p-1"
            style={
              this.props.location.pathname === "/download" ||
              this.props.location.pathname === "/settings"
                ? { visibility: "hidden" }
                : {}
            }
            onClick={() => {
              this.props.navigate("/settings")
            }}
          >
            <Settings fontSize="large" />
          </div>
          <div
            data-bs-toggle="tooltip"
            title="GITHUB"
            data-bs-placement="bottom"
            className="menu-icon p-1"
          >
            <GitHub
              style={{ fontSize: "1.9rem" }}
              onClick={() => {
                if (!window.go) return
                window.go.main.App.OpenGithub()
              }}
            />
          </div>
        </div>
      </header>
    )
  }
}

export default (props) => {
  return (
    <Header
      {...props}
      navigate={useNavigate()}
      location={useLocation()}
      {...useContext(AppContext)}
    />
  )
}
