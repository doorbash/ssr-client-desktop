import React from "react"
import "./DownloadModal.css"

export default function DownloadModal({ title, text, positive, negative }) {
  return (
    <>
      <div className="modal flex justify-center items-center fixed top-0 left-0 w-full h-full outline-none overflow-x-hidden overflow-y-auto bg-[#00000088]">
        <div className="modal-dialog relative w-auto pointer-events-none">
          <div className="modal-content border-none shadow-lg relative flex flex-col w-full pointer-events-auto bg-white bg-clip-padding rounded-md outline-none text-current">
            <div className="modal-body relative p-4">{text}</div>
            <div className="modal-footer flex flex-shrink-0 flex-wrap space-x-1 items-center justify-end p-4 border-t border-gray-200 rounded-b-md">
              {negative && (
                <button
                  type="button"
                  className="btn-cancel"
                  onClick={negative.onClick}
                >
                  {negative.text}
                </button>
              )}
              <button
                type="button"
                className="btn-primary"
                onClick={positive.onClick}
              >
                {positive.text}
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
