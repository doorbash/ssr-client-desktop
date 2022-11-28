import React, { Component, useContext } from "react"
import { useNavigate } from "react-router-dom"
import AppContext from "../../contexts/AppContext"
import DownloadModal from "../DownloadModal/DownloadModal"

class DownloadClient extends Component {
  constructor(props) {
    super(props)
    props.setConfig("Download SSR-Client", false, "#00000000")
    this.onEvent = this.onEvent.bind(this)

    this.state = {
      status: "ready", // ready, ongoing, error, canceled, success
      progress: 0,
      error: "",
    }
  }

  onEvent(key, value, err) {
    console.log("download report >>>> " + key, value, err)
    if (this.state.status === "canceled") return

    if (key === "progress") {
      // update progress bar
      this.setState({
        status: "ongoing",
        progress: value,
      })
    } else if (key === "result") {
      if (value === "canceled") {
        this.setState({
          status: "canceled",
        })
      } else if (value === "error") {
        this.setState({
          status: "error",
          error: err,
        })
      } else if (value === "success") {
        this.setState({
          status: "success",
        })
      }
    }
  }

  componentDidMount() {
    if (!window.runtime) return
    window.runtime.EventsOn("download-report", this.onEvent)
  }

  componentWillUnmount() {
    if (!window.runtime) return
    window.runtime.EventsOff("download-report")
  }

  render() {
    if (this.state.status === "ready") {
      return (
        <DownloadModal
          text="You need to download ssr-client binary file."
          positive={{
            text: "Download",
            onClick: (e) => {
              if (!window.go) return
              this.setState({
                status: "ongoing",
                progress: 0,
                error: "",
              })
              window.go.main.App.DownloadClientFile()
            },
          }}
          negative={{
            text: "Exit",
            onClick: (e) => {
              if (!window.runtime) return
              window.runtime.Quit()
            },
          }}
        />
      )
    } else if (this.state.status === "error") {
      return (
        <DownloadModal
          text={`Error while downloading ssr-client binary file: ${this.state.error}`}
          positive={{
            text: "Try again",
            onClick: (e) => {
              if (!window.go) return
              this.setState({
                status: "ongoing",
                progress: 0,
                error: "",
              })
              window.go.main.App.DownloadClientFile()
            },
          }}
          negative={{
            text: "Exit",
            onClick: (e) => {
              if (!window.runtime) return
              window.runtime.Quit()
            },
          }}
        />
      )
    } else if (this.state.status === "success") {
      return (
        <DownloadModal
          text="Download finished successfully!"
          positive={{
            text: "OK",
            onClick: (e) => {
              this.props.navigate(-1)
            },
          }}
        />
      )
    } else if (this.state.status === "ongoing") {
      return (
        <>
          <div className="p-4">
            <div className="mb-2">Downloading ssr-client...</div>
            <div className="w-full bg-gray-200 mb-4">
              <div
                className={`bg-blue-600  text-xs font-medium text-blue-100 text-center p-0.5 leading-none`}
                style={{ width: `${this.state.progress}%` }}
              >
                <div
                  style={{
                    visibility: `${
                      this.state.progress < 4 ? "hidden" : "unset"
                    }`,
                  }}
                >
                  {this.state.progress}%
                </div>
              </div>
            </div>
            <div className="flex justify-end items-center">
              <button
                className="btn-primary"
                onClick={() => {
                  this.setState({
                    status: "canceled",
                  })

                  if (window.go) {
                    window.go.main.App.CancelDownload()
                  }

                  if (window.runtime) {
                    setTimeout(() => {
                      window.runtime.Quit()
                    }, 500)
                  }
                }}
              >
                Cancel
              </button>
            </div>
          </div>
        </>
      )
    } else {
      return <></>
    }
  }
}

export default (props) => {
  return (
    <DownloadClient
      {...props}
      navigate={useNavigate()}
      {...useContext(AppContext)}
    />
  )
}
