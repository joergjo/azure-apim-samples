using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Mvc.ApiExplorer;

namespace SimpleCalculator
{
    public static class CalculatorConventions
    {
        [ApiConventionNameMatch(ApiConventionNameMatchBehavior.Prefix)]
        [ProducesDefaultResponseType(typeof(string))]
        [ProducesResponseType(200)]
        public static void Get() { }
    }
}
