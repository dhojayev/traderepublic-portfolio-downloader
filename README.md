# Trade Republic Portfolio Downloader

## Preamble

This project exists solely because I need a better representation of all purchases and sales made in Trade Republic.

The initial idea was to create a list of all purchase and sale transactions to be able to track profit and loss.
Unfortunately Trade Republic does not provide such dashboard or file download (it is possible to view current portfolio
assets but no analytics of what has been sold and how much of profit has been gained).

## Existing solutions

After researching a bit in internet I found a solution (namely: https://github.com/Zarathustra2/TradeRepublicApi) (or
two) that fulfill similar requirements, however, those
solutions are pretty outdated and do not benefit from new endpoints (actually websocket message types, but let's call
them here that way for simplicity).

Main disadvantages of existing solutions for me:

* Resets the paired device (which is annoying if you use the mobile app regularly)
* Requiring OCR for getting transaction details
* Being pretty outdated
* Written in python (which I am not proficient in enough to contribute)

## Implementation

### Minimum requirements and limitations

I initially set a few requirements and limitations for myself when planning this project:

* Written in go and compiled as a binary for all major operating systems and platforms
* Making it opensource to allow others to contribute and audit
* No configuration file requirements (all input requested in console)
* No dependencies (such as SQL databases, etc) to enable non-techs using the app "AS IS"
* Writing results into a CSV file to be able to build formulas, filter, etc. the data
* No security information storage (except for session and refresh tokens) on host machine
* **No data should leave the host machine**

### Currently supported functionality

* Creating a CSV file with all transaction (except for "interest received" transaction for now). This includes:
  * purchase and sale of ETFs, stocks, cryptocurrency;
  * dividends received from ETFs and stocks;
  * benefits received suchs as round up and save-back;
* Inserting new data into the CSV file
* Saving raw responses onto the file system

### Planned features and improvements

**What is coming:**

* Support of "interest received" transactions
* Identifying stock transactions
* Downloading and storing PDF files attached to each transaction

**What (maybe) will follow:**

* Writing the data into an SQL database for using it in custom applications

## Usage

### Running the app

After downloading a binary for respective OS simply run it in terminal by providing your phone number as an argument.

```shell
➜  traderepublic-portfolio-downloader git:(main) ✗ ./bin/portfoliodownloader --help
Usage: portfoliodownloader [--write-responses] [--debug] [--trace] PHONENUMBER

Positional arguments:
  PHONENUMBER            Phone number in international format: +49xxxxxxxxxxxxx

Options:
  --write-responses, -w
                         Write api responses to file system
  --debug                Enable debug mode
  --trace                Enable trace mode
  --help, -h             display this help and exit

```

### CSV file fields

| Field               | Description                                                               |
|---------------------|---------------------------------------------------------------------------|
| **ID**              | Transaction UUID                                                          |
| **Status**          | Transaction status (should always be `executed`)                          |
| **Timestamp**       | Date and time of transaction execution, e.g.: `30 Nov 23 10:22 +0000`     |
| **Type**            | Transaction type, one of: `Purchase, Sale, Dividends, Round Up, Saveback` |
| **Asset type**      | Asset type, one of: `ETF, Cryptocurrency, Lending, Other`                 |
| **Name**            | Asset name, e.g.: `Bitcoin`                                               |
| **Instrument**      | Instrument ISIN, e.g.: `IE00BK1PV551`                                     |
| **Shares**          | Number of shares in transaction                                           |
| **Rate**            | Price per share in EUR                                                    |
| **Realized yield**  | Realized yield in percentage                                              |
| **Realized PnL**    | Realized profit or loss amount in EUR (negative is loss)                  |
| **Commission**      | Commission paid to Trade Republic for the transaction in EUR              |
| **Debit**           | Amount debited from the deposited amount in EUR                           |
| **Credit**          | Amount credited to the deposited amount in EUR                            |
| **Portfolio value** | Amount that contributes to the portfolio size in EUR                      |

Example CSV output can be view here: [transactions.csv](assets%2Ftransactions.csv)

## Have suggestions or improvements?

I hope that this app will become better with the help of the community until Trade Republic decides to implement a
better dashboard to be able to have an overview of such data.

Please create a pull request with your changes if you have something to contribute. I am very open for constructive
suggestions and feedback.

## Closing words

This project and I have no affiliation to Trade Republic Bank GmbH by any means. Trade Republic is a registered
trademark of Trade Republic Bank GmbH.

It is important to mention that this application does nothing more than Trade Republic's frontend application would do:

* Authenticates using the same endpoints
* Retrieves the data using the same websocket address
