using System;
using System.Text;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Logging;

namespace SimpleCalculator.Controllers
{
    [ApiController]
    [ApiConventionType(typeof(CalculatorConventions))]
    
    public class CalculatorController : ControllerBase
    {
        private readonly ILogger _logger;
        private static bool _failHealthChecks = false;

        public CalculatorController(ILogger<CalculatorController> logger)
        {
            _logger = logger;
        }

        [Route("api/add")]
        [Produces("application/xml")]
        [HttpGet]
        public ActionResult<string> GetSum([FromQuery] int a, [FromQuery] int b) => CreateResponse(a + b);


        [Route("api/sub")]
        [Produces("application/xml")]
        [HttpGet]
        public ActionResult<string> GetDiff([FromQuery] int a, [FromQuery] int b) => CreateResponse(a - b);

        [Route("api/mul")]
        [Produces("application/xml")]
        [HttpGet]
        public ActionResult<string> GetProduct([FromQuery] int a, [FromQuery] int b) => CreateResponse(a * b);

        [Route("api/div")]
        [Produces("application/xml")]
        [HttpGet]
        public ActionResult<string> GetDivision([FromQuery] int a, [FromQuery] int b) => CreateResponse(a / b);

        [Route("health")]
        [Produces("text/plain")]
        [HttpGet]
        public ActionResult<string> GetHealth()
        {
            if (_failHealthChecks)
            {
                return StatusCode(StatusCodes.Status503ServiceUnavailable);
            }
            return Ok($"Service healthy on '{Environment.MachineName}'.");
        }

        [Route("fail/{active}")]
        [Produces("text/plain")]
        [HttpGet]
        public ActionResult<string> ToggleFailHealthChecks(bool active)
        {
            _failHealthChecks = active;
            string message = $"Service now {(active ? "unhealthy" : "healthy")} on '{Environment.MachineName}'.";
            _logger.LogInformation(message);
            return Ok(message);
        }

        private ActionResult<string> CreateResponse(int result)
        {
            _logger.LogInformation("Calculator result: {Result}.", result);
            string xml = $"<result><value>{result}</value><broughtToYouBy>Azure API Management - http://azure.microsoft.com/apim/ 2020</broughtToYouBy></result>";
            return Content(xml, "application/xml", Encoding.UTF8);
        }
    }
}
