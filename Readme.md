TODO/Ideas:
    1. Use address of store location for buyer and driver.
    2. Driver : Sort by passenger with least detour.
    3. Driver/Buyer : Start at Home address only.

expected response once driver enters destination address:
{
    user[name, address],
    [
        {
            dropoffAddress,
            detourTimeDiff,
            whatToPickup,
            priceRequested
        }
    ]
}

type User {
    email
    name
    phoneNumber
    address
}