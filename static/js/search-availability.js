function searchAvailability(roomId) {
  document.getElementById("check-availability-button").addEventListener("click", function () {
    const html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
          <div class="form-row">
            <div class="col">
              <div class="form-row" id="reservation-dates-modal">
                <div class="col">
                  <input class="form-control" type="text" name="start" id="start" placeholder="Arrival" required disabled />
                </div>
                <div class="col">
                  <input class="form-control" type="text" name="end" id="end" placeholder="Departure" required disabled />
                </div>
              </div>
            </div>
          </div>
        </form>
      `;

    attention.custom({
      msg: html,
      title: "Coose your dates",
      willOpen: () => {
            const elem = document.getElementById('reservation-dates-modal');
            const rp = new DateRangePicker(elem, {
              format: 'yyyy-mm-dd',
              showOnFocus: true,
              minDate: new Date(),
            });
          },
      didOpen: () => {
        document.getElementById('start').removeAttribute('disabled');
        document.getElementById('end').removeAttribute('disabled');
      },
      callback: function(result) {
        const form = document.getElementById("check-availability-form");
        const formData = new FormData(form);
        formData.append("csrf_token", "{{.CSRFToken}}");
        formData.append("room_id", roomId);

        fetch("/search-availability-json", { method: "post", body: formData })
          .then(response => response.json)
          .then(data => {
            if (data.ok) {
              attention.custom({
                showConfirmButton: false,
                showCancelButton: false,
                icon: "success",
                msg: `
                  <p>Room is available!</p>
                  <p>
                    <a
                      href="/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}"
                      class="btn btn-primary"
                    >
                      Book now!
                    </a>
                  </p>
                `,
              });
            } else {
              attention.error({
                msg: "No availability",
              });
            }
          });
      }
     });
  })
}