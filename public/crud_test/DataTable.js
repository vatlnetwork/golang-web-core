import Component from "../js-component-lib/lib/Component.js";
import { useTheme } from "../js-component-lib/lib/useTheme.js";

class DataTable extends Component {
  constructor(props) {
    const { data } = props;
    if (data == undefined || data == null) {
      throw new Error("data must be defined");
    }
    if (!Array.isArray(data)) {
      throw new Error("data must be an array");
    }

    super(document.createElement("table"));

    this.element.style.borderSpacing = "0px";

    this.buildHeader();
    this.initTheme();
  }

  initTheme() {
    useTheme(
      () => {
        this.element.style.color = "white";
      },
      () => {
        this.element.style.color = "black";
      }
    );
  }

  buildHeader() {
    const header = document.createElement("thead");
    const headerRow = document.createElement("tr");
    ["ID", "Number", "Boolean", "Edit", "Delete"].forEach((col) => {
      const cell = document.createElement("th");
      cell.innerHTML = col;
      cell.style.border = "1px solid";
      cell.style.padding = "5px";
      headerRow.appendChild(cell);
    });
    header.appendChild(headerRow);
    this.element.appendChild(header);
  }
}

export default DataTable;
