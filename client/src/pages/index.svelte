<div class="container">
    <div class="column">
        <img class="column-item" src="https://via.placeholder.com/150" alt="asd" />
        <div class="column-item">
            <p>cdn</p>
        </div>
    </div>
    <div class="column">
        <div class="column-item">
            <p>easy</p>
        </div>
        <img class="column-item" src="https://via.placeholder.com/150" alt="asd" />
    </div>
</div>

<style>
    p {
        color: white;
    }

    .container {
        padding-top: 2em;
        overflow: scroll;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    .column {
        padding: 25px;
        display: flex;
        flex-direction: row;
    }

    .column-item {
        padding: 0 10px 0 10px;
    }
</style>
