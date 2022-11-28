import "./App.css"

import Header from "./components/Header/Header"
import ProxyList from "./components/ProxyList/ProxyList"
import NewProxy from "./components/NewProxy/NewProxy"
import { Routes, Route, BrowserRouter as Router } from "react-router-dom"
import DownloadClient from "./components/DownloadClient/DownloadClient"
import Settings from "./components/Settings/Settings"
import Logs from "./components/Logs/Logs"

import React, { Component } from "react"
import AppContext from "./contexts/AppContext"

export default class App extends Component {
  constructor(props) {
    super(props)
    this.setConfig = this.setConfig.bind(this)
    this.updateTooltips = this.updateTooltips.bind(this)
    this.state = {
      title: "",
      showBackButton: false,
      bgColor: "white",
    }
  }

  setConfig(title, showBackButton, bgColor) {
    this.setState({
      title: title,
      showBackButton: showBackButton,
      bgColor: bgColor,
    })
  }

  updateTooltips() {
    if (this.tooltipList) {
      this.tooltipList.forEach((tooltip) => {
        tooltip.dispose()
      })
    }

    this.tooltipList = [].slice
      .call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
      .map(function (tooltipTriggerEl) {
        return new window.Tooltip(tooltipTriggerEl)
      })

    this.tooltipList.forEach((tooltip) => {
      tooltip._config = {
        ...tooltip._config,
        fallbackPlacements: ["bottom"],
        offset: "0,15",
      }
      tooltip
        .getTipElement()
        .querySelector(".tooltip-inner")
        .classList.add("custom-tooltip")
    })
  }

  componentWillUnmount() {
    if (this.tooltipList) {
      this.tooltipList.forEach((tooltip) => {
        tooltip.dispose()
      })
    }
  }

  render() {
    return (
      <div className="app">
        <Router>
          <AppContext.Provider
            value={{
              ...this.state,
              setConfig: this.setConfig,
              updateTooltips: this.updateTooltips,
            }}
          >
            <Header />
            <div
              className="app-content"
              style={{ backgroundColor: this.state.bgColor }}
            >
              <Routes>
                <Route path="/" element={<ProxyList />} />
                <Route path="/new-proxy" element={<NewProxy />} />
                <Route path="/edit-proxy" element={<NewProxy />} />
                <Route path="/download" element={<DownloadClient />} />
                <Route path="/settings" element={<Settings />} />
                <Route path="/logs" element={<Logs />} />
              </Routes>
            </div>
          </AppContext.Provider>
        </Router>
      </div>
    )
  }
}
