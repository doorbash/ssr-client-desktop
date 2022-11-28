import { Add, PlayArrow, Stop, Delete, Edit } from "@material-ui/icons"
import { Component, useContext } from "react"
import { useNavigate } from "react-router-dom"
import "./ProxyList.css"
import "tw-elements"
import AppContext from "../../contexts/AppContext"

class ProxyList extends Component {
  constructor(props) {
    super(props)
    props.setConfig("Shadowsocksr Client", false, "#00000000")
    this.onEvent = this.onEvent.bind(this)
    this.state = {
      proxies: [],
    }
  }

  onEvent(id, status) {
    console.log("proxy status >>>> " + id, status)
    this.setState((st) => {
      return {
        proxies: st.proxies.map((it) => {
          if (it.id === id) {
            return {
              ...it,
              run_status: status, // idle, running, error
            }
          }
          return it
        }),
      }
    })
  }

  componentDidMount() {
    if (!window.go) return
    ;(async () => {
      const list = await window.go.main.App.GetProxies()
      console.log(list)

      this.setState({
        proxies: list,
      })

      const result = await window.go.main.App.ClientFileExists()
      if (result === false) {
        this.props.navigate("/download")
        return
      }

      if (!window.runtime) return
      window.runtime.EventsOn("run-status", this.onEvent)
    })()
  }

  componentWillUnmount() {
    if (!window.runtime) return
    window.runtime.EventsOff("run-status")
  }

  componentDidUpdate() {
    this.props.updateTooltips()
  }

  render() {
    return (
      <>
        {this.state.proxies.map((item, i) => {
          return (
            <div
              key={item.id}
              className="proxy-item cursor-pointer border-b-2 text-lg grid items-center justify-center pl-1"
              onClick={(e) => {
                this.props.navigate("/logs", {
                  state: {
                    id: item.id,
                    name: item.name,
                  },
                })
              }}
            >
              <div
                className="scale-50"
                style={{
                  fill:
                    item.status === 0 || item.run_status !== "running"
                      ? "#c2c2c2"
                      : "#54d1b0",
                }}
              >
                <svg style={{ transform: "scale(-1,1)" }} viewBox="0 0 24 24">
                  <path d="M24 0l-6 22-8.129-7.239 7.802-8.234-10.458 7.227-7.215-1.754 24-12zm-15 16.668v7.332l3.258-4.431-3.258-2.901z" />
                </svg>
              </div>
              <div className="my-4 w-[90%]">
                <div className="mb-[0.125rem] truncate">{item.name}</div>
                <div className="flex justify-start space-x-2 items-center text-xs text-slate-600">
                  {item.status === 0 ? (
                    <div>DISABLED</div>
                  ) : item.run_status === "running" ? (
                    <div>RUNNING</div>
                  ) : item.run_status === "error" ? (
                    <div>ERROR</div>
                  ) : (
                    <div>IDLE</div>
                  )}
                  <div>PORT: {item.l}</div>
                </div>
              </div>
              <div className="proxy-icons">
                <div
                  data-bs-toggle="tooltip"
                  data-bs-placement="bottom"
                  title={item.status === 0 ? "RUN" : "STOP"}
                  className="proxy-icon"
                  onClick={(e) => {
                    e.stopPropagation()
                    if (!window.go) return
                    ;(async () => {
                      try {
                        if (item.status === 0) {
                          await window.go.main.App.RunProxy(item.id)
                        } else {
                          await window.go.main.App.StopProxy(item.id)
                        }
                        this.setState((st) => {
                          return {
                            proxies: st.proxies.map((it) => {
                              if (it.id === item.id) {
                                return {
                                  ...it,
                                  status: 1 - it.status,
                                }
                              }
                              return it
                            }),
                          }
                        })
                      } catch (err) {
                        console.log(">>> ", err)
                      }
                    })()
                  }}
                >
                  {item.status === 0 ? (
                    <PlayArrow style={{ fontSize: "1.3rem" }} />
                  ) : (
                    <Stop style={{ fontSize: "1.3rem" }} />
                  )}
                </div>
                <div
                  data-bs-toggle="tooltip"
                  data-bs-placement="bottom"
                  title="DUPLICATE"
                  className="proxy-icon"
                  onClick={(e) => {
                    e.stopPropagation()
                    this.props.navigate("/new-proxy", {
                      state: {
                        ...item,
                        name: "Copy of " + item.name,
                      },
                    })
                  }}
                >
                  <Add style={{ fontSize: "1.6rem" }} />
                </div>
                <div
                  data-bs-toggle="tooltip"
                  data-bs-placement="bottom"
                  title="EDIT"
                  className="proxy-icon"
                  onClick={(e) => {
                    e.stopPropagation()
                    this.props.navigate("/edit-proxy", { state: item })
                  }}
                >
                  <Edit style={{ fontSize: "1.3rem" }} />
                </div>
                <div
                  data-bs-toggle="tooltip"
                  data-bs-placement="bottom"
                  title="DELETE"
                  offset={1000}
                  className="proxy-icon"
                  onClick={(e) => {
                    e.stopPropagation()
                    if (!window.go) return
                    ;(async () => {
                      try {
                        await window.go.main.App.DeleteProxy(item.id)
                        this.setState((st) => {
                          return {
                            proxies: st.proxies.filter((it) => {
                              return it.id !== item.id
                            }),
                          }
                        })
                      } catch (e) {
                        console.log(e)
                      }
                    })()
                  }}
                >
                  <Delete style={{ fontSize: "1.3rem" }} />
                </div>
              </div>
            </div>
          )
        })}
        <div
          data-bs-toggle="tooltip"
          title="NEW PROXY"
          data-bs-placement="left"
          className="btn-floating absolute right-8 bottom-8"
          onClick={(e) => {
            e.stopPropagation()
            this.props.navigate("/new-proxy")
          }}
        >
          <Add style={{ fontSize: "3rem" }} />
        </div>
      </>
    )
  }
}

export default (props) => {
  return (
    <ProxyList
      {...props}
      navigate={useNavigate()}
      {...useContext(AppContext)}
    />
  )
}
