import Component from "../js-component-lib/lib/Component.js";
import { useTheme } from "../js-component-lib/lib/useTheme.js";

class DataTable extends Component {
  handleDelete;
  data;

  /**
   *
   * @param {{
   *  data: {id: string; number: number; boolean: boolean}[];
   *  handleDelete: (id: string) => {};
   * }} props
   */
  constructor(props) {
    const { data, handleDelete } = props;
    if (data == undefined || data == null) {
      throw new Error("data must be defined");
    }
    if (!Array.isArray(data)) {
      throw new Error("data must be an array");
    }

    super(document.createElement("table"));

    this.element.style.borderSpacing = "0px";

    this.data = data;
    this.handleDelete = handleDelete;

    this.initTheme();
    this.rebuild(data);
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

  /**
   *
   * @param {{id: string; number: number; boolean: boolean}[]} data
   */
  rebuild(data) {
    this.data = data;
    this.element.innerHTML = "";
    this.buildHeader();
    this.buildTable(data);
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

  /**
   *
   * @param {HTMLTableCellElement} cell
   */
  applyCellStyles(cell) {
    cell.style.border = "1px solid";
    cell.style.padding = "5px";
    cell.style.textAlign = "center";
  }

  /**
   *
   * @param {{id: string; number: number; boolean: boolean}[]} data
   */
  buildTable(data) {
    const tbody = document.createElement("tbody");
    data.forEach((item) => {
      const tr = this.buildRow(item);
      tbody.appendChild(tr);
    });
    this.element.appendChild(tbody);
  }

  /**
   *
   * @param {{id: string; number: number; boolean: boolean}} item
   */
  buildRow(item) {
    const tr = document.createElement("tr");
    tr.id = item.id;
    const idCell = this.buildIdCell(item.id);
    const numberCell = this.buildNumberCell(item.number);
    const booleanCell = this.buildBooleanCell(item.boolean);
    const editCell = this.buildEditCell();
    const deleteCell = this.buildDeleteCell(item.id);
    tr.appendChild(idCell);
    tr.appendChild(numberCell);
    tr.appendChild(booleanCell);
    tr.appendChild(editCell);
    tr.appendChild(deleteCell);
    return tr;
  }

  /**
   *
   * @param {string} id
   * @returns
   */
  buildIdCell(id) {
    const idCell = document.createElement("td");
    idCell.innerHTML = id;
    this.applyCellStyles(idCell);
    return idCell;
  }

  /**
   *
   * @param {number} number
   * @returns
   */
  buildNumberCell(number) {
    const numberCell = document.createElement("td");
    numberCell.innerHTML = number.toString();
    this.applyCellStyles(numberCell);
    return numberCell;
  }

  /**
   *
   * @param {boolean} boolean
   * @returns
   */
  buildBooleanCell(boolean) {
    const booleanCell = document.createElement("td");
    booleanCell.innerHTML = boolean.toString();
    this.applyCellStyles(booleanCell);
    return booleanCell;
  }

  buildEditCell() {
    const editCell = document.createElement("td");
    const editButton = document.createElement("button");
    editButton.innerHTML = "Edit";
    editCell.appendChild(editButton);
    this.applyCellStyles(editCell);
    return editCell;
  }

  /**
   *
   * @param {string} id
   * @returns
   */
  buildDeleteCell(id) {
    const deleteCell = document.createElement("td");

    const deleteButton = document.createElement("button");
    deleteButton.innerHTML = "Delete";
    deleteButton.onclick = () => this.handleDelete(id);
    deleteCell.appendChild(deleteButton);

    this.applyCellStyles(deleteCell);

    return deleteCell;
  }
}

export default DataTable;
